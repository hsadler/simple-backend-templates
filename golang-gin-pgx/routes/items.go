package routes

import (
	"context"
	"errors"
	"example-server/dependencies"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ITEMS API

func SetupItemsAPIRoutes(router *gin.Engine, deps *dependencies.Dependencies) {
	itemsRouterGroup := router.Group("/api/items")
	itemsRouterGroup.GET("/:id", HandleGetItem(deps))
	itemsRouterGroup.GET("", HandleGetItems(deps))
	itemsRouterGroup.POST("", HandleCreateItem(deps))
}

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
// @Summary Get Item
// @Description Returns Item by id.
// @Tags items
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} routes.GetItemResponse
// @Failure 400 {object} string
// @Router /api/items/{id} [get]
func HandleGetItem(deps *dependencies.Dependencies) gin.HandlerFunc {
	return func(g *gin.Context) {
		// Parse Item ID
		itemId, err := strconv.Atoi(g.Param("id"))
		if err != nil {
			log.Println("Error parsing Item ID:", err)
			g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item ID"})
			return
		}
		// Fetch Item by ID
		var item Item
		fetchErr := deps.DBPool.QueryRow(
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
}

type GetItemsResponse struct {
	Data []Item   `json:"data"`
	Meta struct{} `json:"meta"`
}

// GetItems godoc
// @Summary Get Items
// @Description Returns Items by ids. Only returns subset of Items found.
// @Tags items
// @Accept json
// @Produce json
// @Param item_ids query []int true "Item IDs" collectionFormat(multi)
// @Success 200 {array} routes.GetItemsResponse
// @Failure 400 {object} string
// @Router /api/items [get]
func HandleGetItems(deps *dependencies.Dependencies) gin.HandlerFunc {
	return func(g *gin.Context) {
		// Parse Item IDs
		var itemIds []int
		var err error
		itemIdsStrArr, ok := g.GetQueryArray("item_ids")
		if ok {
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
			rows, err = deps.DBPool.Query(
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
// @Description Creates Item.
// @Tags items
// @Accept json
// @Produce json
// @Param createItemRequest body routes.CreateItemRequest true "Create Item Request"
// @Success 200 {object} routes.CreateItemResponse
// @Failure 400 {object} string
// @Router /api/items [post]
func HandleCreateItem(deps *dependencies.Dependencies) gin.HandlerFunc {
	return func(g *gin.Context) {
		// Deserialize request
		var createItemRequest CreateItemRequest
		if err := g.ShouldBindJSON(&createItemRequest); err != nil {
			log.Println("Error deserializing request:", err)
			g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}
		// Validate request ItemIn data
		if err := deps.Validator.Struct(createItemRequest.Data); err != nil {
			log.Println("Error validating request:", err)
			g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item data payload"})
			return
		}
		// Insert Item
		var itemId int
		insertErr := deps.DBPool.QueryRow(
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
		fetchErr := deps.DBPool.QueryRow(
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
}
