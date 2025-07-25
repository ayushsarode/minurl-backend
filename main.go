package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ayushsarode/minurl-backend/handlers"
	"github.com/ayushsarode/minurl-backend/middleware"
	"github.com/ayushsarode/minurl-backend/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env file not found, using system environment variables")
	}

	if err := utils.InitDB(); err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	log.Println("✅ Connected to MongoDB")

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
			"https://minurl-frontend-shas-e7xcsyrzx-ayushsarodes-projects.vercel.app",
			"https://minurl-frontend-shas.vercel.app",
		},
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
