package auth

import "testing"

func TestGenerateToken(t *testing.T) {
	secret := "secret"

	token, err := GenerateToken(secret, 1)
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	if token == "" {
		t.Error("expected token to be not empty")
	}
}
