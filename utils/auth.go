package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword : hashing the password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// ComparePassword : compare the password
func ComparePassword(userPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))

	return err == nil
}
