package utils

import (
	"testing"
)

func TestValidateJWT(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDE5LTA0LTIyVDIyOjQ3OjAxLjA1NTU5MSswOTowMCIsInVzZXJJRCI6IjEifQ.3dOqg_ceqv25nCl2C2WL_lye6vLC5dE8nplk184-5lQ"
	userID, err := ValidateJWT(&tokenString)
	if err != nil {
		t.Error(err)
	}

	if userID == nil {
		t.Errorf("Token is expired or something wrong")
	}

	t.Log(userID)
}
