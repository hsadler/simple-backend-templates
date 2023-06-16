package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Example Server API
// @description This is an example server API
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()

	// @Router /status [get]
	// @Summary Get server status
	// @Description Get the server status
	// @Produce json
	// @Success 200 {object} statusResponse
	r.GET("/status", func(c *gin.Context) {
		status := statusResponse{
			Status: "ok!",
		}
		c.JSON(http.StatusOK, status)

	})

	// Register Swagger routes
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The API endpoint URL
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.Run(":8000")
}

// Status response struct
type statusResponse struct {
	Status string `json:"status"`
}

// package main

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/jackc/pgx/v4"
// )

// type Item struct {
// 	ID    int    `json:"id"`
// 	Name  string `json:"name"`
// 	Price int    `json:"price"`
// }

// var db *pgx.Conn

// func main() {
// 	var err error
// 	db, err = pgx.Connect(context.Background(), "postgres://username:password@localhost/mydatabase?sslmode=disable")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	router := gin.Default()

// 	router.POST("/items", createItemHandler)
// 	router.GET("/items/:id", getItemHandler)

// 	log.Fatal(router.Run(":8080"))
// }

// func createItemHandler(c *gin.Context) {
// 	var item Item
// 	if err := c.ShouldBindJSON(&item); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
// 		return
// 	}

// 	err := db.QueryRow(context.Background(), "INSERT INTO items (name, price) VALUES ($1, $2) RETURNING id", item.Name, item.Price).Scan(&item.ID)
// 	if err != nil {
// 		log.Println("Error inserting item:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert item"})
// 		return
// 	}

// 	c.Status(http.StatusCreated)
// }

// func getItemHandler(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
// 		return
// 	}

// 	var item Item
// 	err = db.QueryRow(context.Background(), "SELECT id, name, price FROM items WHERE id = $1", id).Scan(&item.ID, &item.Name, &item.Price)
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
// 		} else {
// 			log.Println("Error retrieving item:", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve item"})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, item)
// }
