package handlers

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	// "fmt"
	"minurl/models"
	"minurl/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Shorten URL
func ShortenURL(c *gin.Context) {
	var url models.URL

	if err := c.ShouldBindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if url.Original == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Original URL is required"})
		return
	}

	collection := utils.GetCollection("urls")

	// Ensure unique short URL
	var existing models.URL
	for {
		url.Short = generateULID()
		err := collection.FindOne(c, bson.M{"short": url.Short}).Decode(&existing)
		if err == mongo.ErrNoDocuments {
			break // Unique short URL found
		}
	}

	url.UserID = c.GetString("userID")
	url.Clicks = 0
    url.CreatedAt = time.Now()

	_, err := collection.InsertOne(c, url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create URL"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"original_url": url.Original,
		"short_url":    url.Short,
	})
}

// Generate a short URL using ULID + Base64
func generateULID() string {
	t := time.Now()
	entropy := rand.Reader
	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	shortID := base64.RawURLEncoding.EncodeToString(id.Bytes())
	return strings.ToLower(shortID[:6]) // Reduce length to prevent collisions
}

// Redirect URL
func RedirectURL(c *gin.Context) {
	short := c.Param("short")
	var url models.URL

	collection := utils.GetCollection("urls")
	err := collection.FindOne(c, bson.M{"short": short}).Decode(&url)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Update click count directly using ID (No need to convert ObjectID again)
	_, err = collection.UpdateOne(c, bson.M{"_id": url.ID}, bson.M{"$inc": bson.M{"clicks": 1}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update click count"})
		return
	}

	// CORS Headers
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

	// Redirect to the original URL
	c.Redirect(http.StatusFound, url.Original)
}

// Get All User Links
func GetUserLinks(c *gin.Context) {
    collection := utils.GetCollection("urls")
    userID := c.GetString("userID")
    fmt.Printf("UserID from context: %s\n", userID)
    
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    var urls []models.URL
    // Changed "user_id" to "userID" to match the model's bson tag
    cursor, err := collection.Find(context.TODO(), bson.M{"userID": userID}, options.Find().SetSort(bson.M{"createdAt": -1}))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URLs"})
        return
    }
    defer cursor.Close(context.TODO())
    
    if err := cursor.All(context.TODO(), &urls); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse URLs"})
        return
    }
    
    if urls == nil { // Ensure we return an array, not null
        urls = []models.URL{}
    }
    
    c.JSON(http.StatusOK, gin.H{
        "urls": urls,
        "count": len(urls),
    })
}

func DeleteURL(c *gin.Context) {
    // Get URL ID from parameters
    shortCode := c.Param("short")
    if shortCode == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Short code is required"})
        return
    }

    // Get userID from context (set by auth middleware)
    userID := c.GetString("userID")
    if userID == "" {
        fmt.Println("[ERROR] Unauthorized request: Missing userID in context")
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    fmt.Println("[INFO] Deleting URL:", shortCode, "for user:", userID)

    collection := utils.GetCollection("urls")

    // Delete URL ensuring it belongs to the user
    result, err := collection.DeleteOne(context.TODO(), bson.M{
        "short":  shortCode,
        "userID": userID,
    })

    if err != nil {
        fmt.Printf("[ERROR] Failed to delete URL: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete URL", "details": err.Error()})
        return
    }

    if result.DeletedCount == 0 {
        fmt.Println("[WARNING] URL not found or unauthorized:", shortCode)
        c.JSON(http.StatusNotFound, gin.H{"error": "URL not found or unauthorized"})
        return
    }

    fmt.Println("[SUCCESS] URL deleted:", shortCode)
    c.JSON(http.StatusOK, gin.H{
        "message": "URL deleted successfully",
    })
}



func GenerateQRCode(c *gin.Context) {
    // Get the short URL code from the request parameters
    shortCode := c.Param("short")
    if shortCode == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Short code is required"})
        return
    }

    // Get the collection
    collection := utils.GetCollection("urls")

    // Find the URL in the database
    var url models.URL
    err := collection.FindOne(c, bson.M{"short": shortCode}).Decode(&url)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
        return
    }

    // Generate the full URL for the QR code
    fullURL := fmt.Sprintf("http://localhost:8080/%s", shortCode) // Replace with your domain

    // Generate QR code
    var qr *qrcode.QRCode
    qr, err = qrcode.New(fullURL, qrcode.Medium)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
        return
    }

    // Create a buffer to store the PNG
    var buf bytes.Buffer
    err = qr.Write(256, &buf) // 256x256 pixels
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code image"})
        return
    }

    // Set content type header
    c.Header("Content-Type", "image/png")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.png\"", shortCode))
    
    // Set CORS headers if needed
    c.Header("Access-Control-Allow-Origin", "http://localhost:5174")
    c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
    c.Header("Access-Control-Allow-Headers", "Origin, Content-Type")

    // Write the image to the response
    c.Data(http.StatusOK, "image/png", buf.Bytes())
}