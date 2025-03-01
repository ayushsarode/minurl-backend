package handlers

import (
	"context"
	"fmt"
	"github.com/ayushsarode/minurl-backend/utils"
	"github.com/ayushsarode/minurl-backend/models"
	"net/http"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	
	user.Password = string(hashedPassword)
	collection := utils.GetCollection("users")
	_, err = collection.InsertOne(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var dbUser models.User
	collection := utils.GetCollection("users")
	err := collection.FindOne(c, gin.H{"email": user.Email}).Decode(&dbUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	
	token, err := utils.GenerateToken(dbUser.ID.Hex()) // Convert ObjectID to string
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	
	// Return the profile picture URL in the login response
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":          dbUser.ID.Hex(),
			"name":        dbUser.Username,
			"email":       dbUser.Email,
			"profile_pic": dbUser.ProfilePic, 
		},
	})
}


func UploadProfilePic(c *gin.Context) {
	// Get user ID from request parameter
	userID := c.Param("id")
	
	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	// Get the file from form-data
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get image"})
		return
	}
	
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer src.Close()
	
	// Initialize Cloudinary
	// Replace these values with your actual Cloudinary credentials
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary"})
		return
	}
	
	// Create a unique public ID for the image
	publicID := fmt.Sprintf("profile_pics/%s_%s", userID, uuid.New().String())
	
	// Upload file to Cloudinary
	uploadResult, err := cld.Upload.Upload(context.Background(), src, uploader.UploadParams{
		PublicID:     publicID,
		ResourceType: "image",
		Folder:       "profile_pics",
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to Cloudinary: " + err.Error()})
		return
	}
	
	// Get the secure URL from the upload result
	profilePicURL := uploadResult.SecureURL
	
	// Update user's profile picture in MongoDB
	collection := utils.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"profile_pic": profilePicURL}}
	
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile picture: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Upload successful", "url": profilePicURL})
}