package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var AllowedOrigins = []string{
	"https://api.maziwasoko.co.ke",
	"http://localhost",
	"http://127.0.0.1",
	"http://192.168.0.103",
	"http://192.168.0.100:4001",
	"http://localhost:4001",
}

var AllowedHeaders = []string{
	"Authorization", "Accept", "Accept-Charset", "Accept-Language",
	"Accept-Encoding", "Origin", "Host", "User-Agent", "Content-Length",
	"Content-Type", "X-Authorization", "Access-Control-Allow-Origin",
	"Access-Control-Allow-Methods", "Access-Control-Allow-Headers",
}

func CORSMiddleware() gin.HandlerFunc {
	// Define and return the CORS middleware directly
	config := cors.Config{
		AllowOrigins:     AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "UPDATE", "OPTIONS"},
		AllowHeaders:     AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	return cors.New(config)
}
