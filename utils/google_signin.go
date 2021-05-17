package utils

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/idtoken"
)

type Claims struct {
	Sub   string
	Email string
}

func GoogleSignIn(ctx context.Context, idToken string) (*Claims, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	payload, err := idtoken.Validate(ctx, idToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return nil, err
	}

	claims := Claims{}
	err = mapstructure.Decode(payload.Claims, &claims)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}
