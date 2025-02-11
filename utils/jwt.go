package utils

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))


func GenerateToken(userID string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,               
		"exp": time.Now().Add(time.Hour * 24).Unix(), 
	})


	return token.SignedString(jwtSecret)
}


func ValidateToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		
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
