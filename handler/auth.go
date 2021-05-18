package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/mattdamon108/go-graphql-api-boilerplate/utils"
)

// ContextKey for the userID in context
type ContextKey string

// Authenticate for JWT
func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// var userID *string

		ctx := r.Context()
		userID := validateAuthHeader(ctx, r)

		if userID != nil {
			ctx = context.WithValue(ctx, ContextKey("userID"), *userID)
		}

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateAuthHeader(ctx context.Context, r *http.Request) *string {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil
	}

	authorization := strings.Split(tokenString, " ")
	if len(authorization) != 2 || authorization[0] != "Bearer" {
		return nil
	}

	userID, _ := utils.ValidateJWT(&authorization[1])
	return userID
}
