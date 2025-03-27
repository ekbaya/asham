package utilities

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse defines the structure of the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Show(c *gin.Context, code int, msg interface{}, data interface{}) {
	c.JSON(code, gin.H{
		"success":     true,
		"status_code": code,
		"message":     msg,
		"data":        data,
	})

}

func ShowMessage(c *gin.Context, code int, msg interface{}) {
	c.JSON(code, gin.H{
		"success":     code == http.StatusOK || code == http.StatusCreated,
		"status_code": code,
		"message":     msg,
	})
}

func ShowError(c *gin.Context, code int, errors []string) {
	c.JSON(code, gin.H{
		"success":     false,
		"status_code": code,
		"error":       errors,
	})
}
