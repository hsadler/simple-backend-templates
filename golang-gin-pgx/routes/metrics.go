package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

// Metrics godoc
// @Summary Metrics
// @Description Returns Prometheus metrics.
// @Tags metrics
// @Produce text/plain
// @Success 200 {string} string
// @Router /metrics [get]
func HandleMetrics(g *gin.Engine) gin.HandlerFunc {
	log.Info().Msg("Request to /metrics")
	return gin.WrapH(promhttp.Handler())
}
