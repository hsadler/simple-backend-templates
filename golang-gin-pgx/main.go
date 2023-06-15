package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Server is running!")
	})

	http.ListenAndServe(":8080", nil)
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
