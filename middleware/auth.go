package middleware

import (
	// "fmt"
	"minurl/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates the JWT token in the Authorization header
func Authmiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Ensure that the token is prefixed with 'Bearer '
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Split the tokenString into parts: should be "Bearer <token>"
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Extract the token without the 'Bearer' prefix
		tokenString = parts[1]

		// Validate the token
		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract userID from claims and pass it to the next handler
		claims := token.Claims.(jwt.MapClaims)
		c.Set("userID", claims["userID"])

		c.Next()
	}
}
