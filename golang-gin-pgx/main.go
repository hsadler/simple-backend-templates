package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerfiles "github.com/swaggo/gin-swagger/swaggerFiles"

	_ "example-server/docs"
)

func main() {
	r := gin.Default()
	r.GET("/status", HandleStatus)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8000")
}

// Status godoc
// @Summary status endpoint
// @Schemes
// @Description returns "ok!" if server is up
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
