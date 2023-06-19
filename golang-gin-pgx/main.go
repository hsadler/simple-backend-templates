package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerfiles "github.com/swaggo/gin-swagger/swaggerFiles"

	_ "example-server/docs"
)

// Global variables
// var db *pgx.Conn
var dbpool *pgxpool.Pool
var validate *validator.Validate

// @title Example Server API
// @description Example Go+Gin+pgx JSON API server.
// @version 1
// @host localhost:8000
// @BasePath /
// @schemes http
func main() {

	// Connect to database and create tables
	var err error
	dbpool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer dbpool.Close()
	log.Println("Connected to database")
	CreateTables()
	log.Println("Created tables")

	// Setup validator
	validate = validator.New()

	// Setup Gin router
	r := gin.Default()

	r.GET("/status", HandleStatus)

	itemsRouterGroup := r.Group("/api/items")
	itemsRouterGroup.GET("/:id", HandleGetItem)
	itemsRouterGroup.GET("", HandleGetItems)
	itemsRouterGroup.POST("", HandleCreateItem)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Run server
	log.Fatal(r.Run(":8000"))
}

func CreateTables() {
	_, err := dbpool.Exec(context.Background(), `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS item (
			id SERIAL PRIMARY KEY,
			uuid UUID DEFAULT uuid_generate_v4(),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			name VARCHAR(50),
			price NUMERIC(10, 2),
			CONSTRAINT name_unique UNIQUE (name)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

// STATUS API

type statusResponse struct {
	Status string `json:"status" example:"ok!"`
}

// Status godoc
// @Summary status endpoint
// @Description Returns `"ok!"` if the server is up
// @Tags status
// @Produce json
// @Success 200 {object} statusResponse
// @Router /status [get]
func HandleStatus(g *gin.Context) {
	status := statusResponse{
		Status: "ok!",
	}
	g.JSON(http.StatusOK, status)
}

// ITEM API

type ItemIn struct {
	Name  string   `json:"name" example:"foo" format:"string" validate:"required"`
	Price *float32 `json:"price" example:"3.14" format:"float64" validate:"min=0"`
}

type Item struct {
	ID        int       `json:"id" example:"1" format:"int64"`
	UUID      string    `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00.000Z" format:"date-time"`
	Name      string    `json:"name" example:"foo" format:"string"`
	Price     float32   `json:"price" example:"3.14" format:"float64"`
}

type GetItemResponse struct {
	Data Item     `json:"data"`
	Meta struct{} `json:"meta"`
}

// GetItem godoc
// @Summary get Item by id
// @Description Returns Item by id
// @Tags items
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} main.GetItemResponse
// @Failure 400 {object} string
// @Router /api/items/{id} [get]
func HandleGetItem(g *gin.Context) {
	// Parse Item ID
	itemId, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		log.Println("Error parsing Item ID:", err)
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item ID"})
		return
	}
	// Fetch Item by ID
	var item Item
	fetchErr := dbpool.QueryRow(
		context.Background(),
		"SELECT id, uuid, created_at, name, price FROM item WHERE id = $1",
		itemId,
	).Scan(&item.ID, &item.UUID, &item.CreatedAt, &item.Name, &item.Price)
	// Handle Item fetch error
	if fetchErr != nil {
		log.Println("Error querying Item:", fetchErr)
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query Item"})
		return
	}
	// Return response
	g.JSON(http.StatusOK, GetItemResponse{Data: item, Meta: struct{}{}})
}

type GetItemsResponse struct {
	Data []Item   `json:"data"`
	Meta struct{} `json:"meta"`
}

// GetItems godoc
// @Summary Get Items
// @Description Returns Items by ids
// @Tags items
// @Accept json
// @Produce json
// @Param item_ids query []int true "Item IDs"
// @Success 200 {array} main.GetItemsResponse
// @Failure 400 {object} string
// @Router /api/items [get]
func HandleGetItems(g *gin.Context) {
	// Parse Item IDs
	var itemIds []int
	var err error
	if itemIdsStr := g.Query("item_ids"); itemIdsStr != "" {
		itemIdsStrArr := strings.Split(itemIdsStr, ",")
		itemIds = make([]int, len(itemIdsStrArr))
		for i, itemIdStr := range itemIdsStrArr {
			itemIds[i], err = strconv.Atoi(itemIdStr)
			// Handle Item ID parse error
			if err != nil {
				log.Println("Error parsing Item ID:", err)
				g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item ID"})
				return
			}
		}
	}
	// Fetch Items by IDs
	var items []Item
	var rows pgx.Rows
	if len(itemIds) > 0 {
		rows, err = dbpool.Query(
			context.Background(),
			"SELECT id, uuid, created_at, name, price FROM item WHERE id = ANY($1)",
			itemIds,
		)
	}
	// Handle Items fetch error
	if err != nil {
		log.Println("Error querying Items:", err)
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query Items"})
		return
	}
	defer rows.Close()
	// Iterate over Items
	for rows.Next() {
		var item Item
		// Scan Item and append to Items unless error
		if err := rows.Scan(&item.ID, &item.UUID, &item.CreatedAt, &item.Name, &item.Price); err != nil {
			log.Println("Error scanning Item:", err)
			g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan Item"})
			return
		}
		items = append(items, item)
	}
	// Return response
	g.JSON(http.StatusOK, GetItemsResponse{Data: items, Meta: struct{}{}})
}

type CreateItemRequest struct {
	Data ItemIn `json:"data"`
}

type CreateItemResponseMeta struct {
	Created bool `json:"created"`
}

type CreateItemResponse struct {
	Data Item                   `json:"data"`
	Meta CreateItemResponseMeta `json:"meta"`
}

// CreateItem godoc
// @Summary Create Item
// @Description Creates Item
// @Tags items
// @Accept json
// @Produce json
// @Param createItemRequest body main.CreateItemRequest true "Create Item Request"
// @Success 200 {object} main.CreateItemResponse
// @Failure 400 {object} string
// @Router /api/items [post]
func HandleCreateItem(g *gin.Context) {
	// Deserialize request
	var createItemRequest CreateItemRequest
	if err := g.ShouldBindJSON(&createItemRequest); err != nil {
		log.Println("Error deserializing request:", err)
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}
	// Validate request ItemIn data
	if err := validate.Struct(createItemRequest.Data); err != nil {
		log.Println("Error validating request:", err)
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item data payload"})
		return
	}
	// Insert Item
	var itemId int
	insertErr := dbpool.QueryRow(
		context.Background(),
		"INSERT INTO item (name, price) VALUES ($1, $2) RETURNING id",
		createItemRequest.Data.Name,
		createItemRequest.Data.Price,
	).Scan(&itemId)
	// Handle Item insert error
	if insertErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(insertErr, &pgErr) {
			// Duplicate entry error handling
			if pgErr.Code == "23505" {
				log.Println("Duplicate Item entry error:", pgErr)
				g.JSON(
					http.StatusConflict,
					gin.H{"error": "Item already exists with name '" + createItemRequest.Data.Name + "'"},
				)
				return
			}
		}
		log.Println("Error inserting Item:", insertErr)
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Item"})
		return
	}
	log.Printf("Inserted itemId: %+v\n", itemId)
	// Fetch Item after insert
	var item Item
	fetchErr := dbpool.QueryRow(
		context.Background(),
		"SELECT id, uuid, created_at, name, price FROM item WHERE id = $1",
		itemId,
	).Scan(&item.ID, &item.UUID, &item.CreatedAt, &item.Name, &item.Price)
	if fetchErr != nil {
		log.Println("Error querying Item:", fetchErr)
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query Item"})
		return
	}
	// Return response
	g.JSON(
		http.StatusOK,
		CreateItemResponse{Data: item, Meta: CreateItemResponseMeta{Created: true}},
	)
}
