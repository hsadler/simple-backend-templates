package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"example-server/database"
	"example-server/dependencies"
	_ "example-server/docs"
	"example-server/logger"
	"example-server/routes"
)

func init() {
	// Setup global logger
	logger.SetupGlobalLogger()
	// Gin settings
	gin.DefaultWriter = os.Stdout
	gin.SetMode(gin.DebugMode)
}

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
	// Swagger docs
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Setup API routes
	routes.SetupItemsAPIRoutes(r, deps)
	// Run server
	log.Info().Msg("Starting server")
	err := r.Run(":8000")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
