package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type statusResponse struct {
	Status string `json:"status" example:"ok"`
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
