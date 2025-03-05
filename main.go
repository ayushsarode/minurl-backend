package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ayushsarode/minurl-backend/handlers"
	"github.com/ayushsarode/minurl-backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, using environment variables")
	}


	// Port configuration
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
		log.Println("PORT not set, defaulting to 8080")
	}

	// Gin setup
	gin.SetMode(gin.ReleaseMode)
	route := gin.Default()

	// CORS configuration
	route.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			os.Getenv("FRONTEND_URL"), // Add your frontend URL from .env
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
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
	route.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Start the server
	log.Printf("Starting server on port %s", httpPort)
	if err := route.Run(":" + httpPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}