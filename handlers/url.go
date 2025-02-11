package handlers

import (
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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	c.Header("Access-Control-Allow-Origin", "http://localhost:5174")
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
    cursor, err := collection.Find(context.TODO(), bson.M{"userID": userID})
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