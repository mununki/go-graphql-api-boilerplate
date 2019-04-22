package utils

import (
	"testing"
)

func TestSignJWT(t *testing.T) {
	userID := "1"
	token, err := SignJWT(&userID)
	if err != nil {
		t.Error(err)
	}
	t.Log(*token)
}
