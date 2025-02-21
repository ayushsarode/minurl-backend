package main

import (
	"minurl/handlers"
	"minurl/middleware"
	// "net/http"
"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)



func main() {
	route := gin.Default()


	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Change this to your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))


	// Public Routes
	route.POST("/register", handlers.Register)
	route.POST("/login", handlers.Login)

	// Protected Routes (Require Auth)
	auth := route.Group("/")
	auth.Use(middleware.Authmiddleware())
	{
		auth.POST("/shorten", handlers.ShortenURL)
		auth.GET("/urls", handlers.GetUserLinks)
		auth.DELETE("/urls/:short", handlers.DeleteURL)
		

	}

	// Public Route for Redirection
	route.GET("/:short", handlers.RedirectURL)
	route.GET("/qr/:short", handlers.GenerateQRCode)
	
	// route.GET("clicks/:short", handlers.GetClickCount)
	

	// Start the server
	route.Run(":8080")
}
