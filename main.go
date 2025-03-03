package main

import (
	"net/http"

	"github.com/ayushsarode/minurl-backend/handlers"
	"github.com/ayushsarode/minurl-backend/middleware"

	// "net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)



func main() {
	route := gin.Default()


	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Change this to your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))


	route.Static("/uploads", "./uploads")


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
		auth.POST("/upload/:id", handlers.UploadProfilePic)
	}

	// Public Route for Redirection
	route.GET("/:short", handlers.RedirectURL)
	route.GET("/qr/:short", handlers.GenerateQRCode)


	route.GET("/ping", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H {
			"message": "pong",
		})
	})

	// Start the server

    httpPort := os.Getenv("PORT")
    if httpPort == "" {
        httpPort = "8000"
    }

	route.Run(":" + httpPort)
}
