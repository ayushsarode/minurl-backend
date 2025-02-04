package handlers

import (
	"minurl/models"
	"minurl/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User

	if err:= c.ShouldBindJSON(&user); err != nil {
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

	err := collection.FindOne(c, gin.H{"username": user.Username}).Decode(&dbUser)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password),[]byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(dbUser.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not ge\nerate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}