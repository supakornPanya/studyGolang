package user_test

import (
	"study-golang-backend/internal/user"
	"testing"
)

func TestMemoryRepository(t *testing.T) {
	repo := user.NewMemoryRepository()

	u := &user.User{
		ID:           "1",
		Username:     "testuser",
		PasswordHash: "hashed_password",
	}

	err := repo.Create(u)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Test duplicate username
	err = repo.Create(u)
	if err == nil {
		t.Fatal("expected error when creating duplicate username, got nil")
	}

	fetched, err := repo.GetByUsername("testuser")
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if fetched.Username != u.Username {
		t.Errorf("expected username %s, got %s", u.Username, fetched.Username)
	}

	_, err = repo.GetByUsername("nonexistent")
	if err == nil {
		t.Fatal("expected error when getting nonexistent user, got nil")
	}
}
