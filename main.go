package main

import (
	"minurl/handlers"
	"minurl/middleware"

	"github.com/gin-gonic/gin"
	// "minurl/middleware"
)

func main() {
	route := gin.Default()

	// Public routes


	// Public routes
	route.POST("/register", handlers.Register)
	route.POST("/login", handlers.Login)

	// Protected routes - AuthMiddleware is used here to check JWT
	auth := route.Group("/")
	auth.Use(middleware.Authmiddleware())
	{
		auth.POST("/shorten", handlers.ShortenURL) // URL shortening requires valid token
	}

	// Redirect route for short URLs
	route.GET("/:short", handlers.RedirectURL)

	// Start the server
	route.Run(":8080")
}
