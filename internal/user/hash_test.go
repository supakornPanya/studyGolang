package user_test

import (
	"study-golang-backend/internal/user"
	"testing"
)

func TestPasswordHashing(t *testing.T) {
	password := "supersecure123"

	hashed, err := user.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if hashed == password {
		t.Fatal("hashed password matches plain password")
	}

	if !user.CheckPasswordHash(password, hashed) {
		t.Fatal("failed to verify correct password hash")
	}

	if user.CheckPasswordHash("wrongpassword", hashed) {
		t.Fatal("verified incorrect password hash")
	}
}
