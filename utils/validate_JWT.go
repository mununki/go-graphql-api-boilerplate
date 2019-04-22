package utils

import (
	// "encoding/base64"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

// ValidateJWT : func to parse JWT and to return the identity
func ValidateJWT(tokenString *string) (*string, error) {
	token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("  Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("my_secret"), nil
	})

	if err != nil {
		return nil, err
	}

	var userID string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID = claims["userID"].(string)
	} else {
		// should do something here!
		return nil, err
	}

	return &userID, err
}
