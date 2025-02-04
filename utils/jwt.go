package utils

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

// Get the JWT secret from the environment variable
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateToken creates and signs a JWT token for a given user ID
func GenerateToken(userID string) (string, error) {
	// Create a new JWT token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,                 // Set userID as a claim
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Set token expiration (1 day)
	})

	// Sign the token with the JWT secret
	return token.SignedString(jwtSecret)
}

// ValidateToken parses and validates the JWT token string
func ValidateToken(tokenString string) (*jwt.Token, error) {
	// Parse the token string with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method is HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
