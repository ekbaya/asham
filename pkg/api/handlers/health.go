package handlers

import (
	"net/http"

	"github.com/ekbaya/asham/pkg/utilities"
	"github.com/gin-gonic/gin"
)

// HealthCheckHandler handles the /health endpoint
func HealthCheckHandler(c *gin.Context) {
	response := utilities.HealthResponse{
		Status:  "OK",
		Message: "Service is running",
	}
	c.JSON(http.StatusOK, response)
}
