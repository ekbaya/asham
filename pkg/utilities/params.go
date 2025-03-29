package utilities

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func IntQueryParam(c *gin.Context, paramName string, defaultValue int) int {
	paramStr := c.Query(paramName)
	if paramStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(paramStr)
	if err != nil {
		// If there's an error parsing the integer, return the default value
		return defaultValue
	}

	return value
}
