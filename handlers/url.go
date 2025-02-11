package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"minurl/models"
	"minurl/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
)

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

	url.Short = generateULID()
	url.UserID = c.GetString("userID")
	url.Clicks = 0

	collection := utils.GetCollection("urls")

	_, err := collection.InsertOne(c, url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create URL"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"original_url": url.Original,
		"short_url":  url.Short,
	})
}

func generateULID() string {
	t := time.Now()
	entropy := rand.Reader

	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	shortID := base64.RawURLEncoding.EncodeToString(id.Bytes())
	
	lowerID := strings.ToLower(shortID)

	return lowerID[:6]
	
}


func RedirectURL(c *gin.Context) {
	short := c.Param("short")

	var url models.URL
	collection := utils.GetCollection("urls")
	err := collection.FindOne(c, bson.M{"short": short}).Decode(&url)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Increment the click count
	updateResult, err := collection.UpdateOne(
		c,
		bson.M{"_id": url.ID},
		bson.M{"$inc": bson.M{"clicks": 2}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update click count", "details": err.Error()})
		return
	}

	// Log the update result for debugging
	fmt.Printf("Update Result: Matched %v, Modified %v\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	c.Redirect(http.StatusMovedPermanently, url.Original)
}

// handlers/url.go
func GetClickCount(c *gin.Context) {
	short := c.Param("short")

	var url models.URL
	collection := utils.GetCollection("urls")
	err := collection.FindOne(c, bson.M{"short": short}).Decode(&url)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clicks": url.Clicks})
}