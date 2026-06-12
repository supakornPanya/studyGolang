package auth_test

import (
	"study-golang-backend/internal/auth"
	"testing"
)

func TestTokenService(t *testing.T) {
	secret := []byte("secret")
	username := "user1"

	tokenString, err := auth.GenerateToken(username, secret)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	parsedUsername, err := auth.ValidateToken(tokenString, secret)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if parsedUsername != username {
		t.Errorf("expected username %s, got %s", username, parsedUsername)
	}

	// Test validation with wrong secret
	_, err = auth.ValidateToken(tokenString, []byte("wrong_secret"))
	if err == nil {
		t.Fatal("expected validation failure with incorrect secret, got nil")
	}
}
