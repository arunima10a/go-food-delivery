package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTestToken( secret string, userID string, role string) ( string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID, 
		"role" : role, 
		"exp" : time.Now().Add(time.Hour * 1).Unix(),
	})
	return token.SignedString([]byte(secret))
}

func ValidateTestToken(tokenString string, secret string) (*jwt.Token, error){
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
		return []byte(secret), nil
	})
}