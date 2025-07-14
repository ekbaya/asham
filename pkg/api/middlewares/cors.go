package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var AllowedOrigins = []string{
	"https://api.maziwasoko.co.ke",
	"https://master.d132eek3o0f1jh.amplifyapp.com",
	"https://arsolibrary.mystock.co.ke",
	"https://www.arsolibrary.mystock.co.ke",
	"http://localhost",
	"http://127.0.0.1:3000",
	"http://192.168.100.43:3000",
	"http://192.168.100.43:5173",
	"http://41.80.113.26:3000",
	"http://192.168.100.2:3000",
	"http://192.168.100.7:3000",
	"http://192.168.100.7:3000",
	"http://localhost:5173",
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
