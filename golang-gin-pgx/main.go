package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"example-server/database"
	"example-server/dependencies"
	_ "example-server/docs"
	"example-server/routes"
)

// @title Example Server API
// @description Example Go+Gin+pgx JSON API server.
// @version 1
// @host localhost:8000
// @BasePath /
// @schemes http
func main() {
	// Setup dependencies
	dbPool, _ := database.SetupDB()
	deps := dependencies.NewDependencies(
		validator.New(),
		dbPool,
	)
	defer deps.CleanupDependencies()
	// Setup Gin router
	r := gin.Default()
	// Status
	r.GET("/status", routes.HandleStatus)
	// Prometheus metrics
	r.GET("/metrics", routes.HandleMetrics(r))
	// Setup API routes
	routes.SetupItemsAPIRoutes(r, deps)
	// Swagger docs
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Run server
	log.Fatal(r.Run(":8000"))
}
