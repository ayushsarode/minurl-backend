package handlers

import (
	"math/rand"
	"minurl/models"
	"minurl/utils"
	"net/http"
	"time"


	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
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

	collection := utils.GetCollection("urls")

	_, err := collection.InsertOne(c, url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create URL"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"original_url": url.Original,
		"short_url": url.Short,
	})
}

func generateULID() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func RedirectURL(c *gin.Context) {
	short := c.Param("short")

	if short == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short URL is required"})
	}

	var url models.URL

	collection := utils.GetCollection("urls")

	err := collection.FindOne(c, gin.H{"short": short}).Decode(&url)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URl not found"})
		return 
	
	}
	c.Redirect(http.StatusMovedPermanently, url.Original)

}