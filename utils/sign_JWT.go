package utils

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// SignJWT : func to generate JWT
func SignJWT(userID *string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": *userID,
		"exp":    time.Now().Add(time.Second * 30 * 24 * 60 * 60),
	})

	tokenString, err := token.SignedString([]byte("my_secret"))

	return &tokenString, err
}
