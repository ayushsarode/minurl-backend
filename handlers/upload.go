package handlers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	// "minurl/models"
	"minurl/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// UploadProfilePic handles profile picture upload and updates MongoDB
func UploadProfilePic(c *gin.Context) {
	// Get user ID from request parameter
	userID := c.Param("id")

	// Get the file from form-data
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get image"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join("uploads", newFilename)

	// Save the file locally
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Update user's profile picture in MongoDB
	profilePicURL := fmt.Sprintf("http://localhost:8080/%s", filePath)
	collection := utils.GetCollection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"profile_pic": profilePicURL}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile picture"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Upload successful", "url": profilePicURL})
}
