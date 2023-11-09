package routes

import (
	"example-server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Status godoc
// @Summary Status
// @Description Returns `"ok"` if the server is up.
// @Tags status
// @Produce json
// @Success 200 {object} models.StatusResponse
// @Router /status [get]
func HandleStatus(g *gin.Context) {
	status := models.StatusResponse{
		Status: "ok",
	}
	g.JSON(http.StatusOK, status)
}
