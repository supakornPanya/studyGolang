package auth_test

import (
	"net/http"
	"net/http/httptest"
	"study-golang-backend/internal/auth"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := []byte("secret")

	r := gin.New()
	r.Use(auth.AuthMiddleware(secret))
	r.GET("/protected", func(c *gin.Context) {
		username, _ := c.Get("username")
		c.String(http.StatusOK, "welcome "+username.(string))
	})

	// 1. Unauthenticated Request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}

	// 2. Authenticated Request
	token, err := auth.GenerateToken("john", secret)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "welcome john" {
		t.Errorf("expected body 'welcome john', got %s", w.Body.String())
	}
}
