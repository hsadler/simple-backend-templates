package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerfiles "github.com/swaggo/gin-swagger/swaggerFiles"

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
	validator := validator.New()
	dbPool := database.SetupDB()
	defer dbPool.Close()
	deps := dependencies.Dependencies{
		Validator: validator,
		DBPool:    dbPool,
	}
	// Setup Gin router
	r := gin.Default()
	// Status
	r.GET("/status", HandleStatus)
	// Prometheus metrics
	r.GET("/metrics", HandleMetrics(r))
	// Setup API routes
	routes.SetupItemsAPIRoutes(r, &deps)
	// Swagger docs
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// Run server
	log.Fatal(r.Run(":8000"))
}

type statusResponse struct {
	Status string `json:"status" example:"ok!"`
}

// Status godoc
// @Summary Status
// @Description Returns `"ok"` if the server is up.
// @Tags status
// @Produce json
// @Success 200 {object} statusResponse
// @Router /status [get]
func HandleStatus(g *gin.Context) {
	status := statusResponse{
		Status: "ok",
	}
	g.JSON(http.StatusOK, status)
}

// Metrics godoc
// @Summary Metrics
// @Description Returns Prometheus metrics.
// @Tags metrics
// @Produce text/plain
// @Success 200 {string} string
// @Router /metrics [get]
func HandleMetrics(g *gin.Engine) gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}
