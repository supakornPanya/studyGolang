package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"study-golang-backend/internal/user"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUserHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := user.NewMemoryRepository()
	h := user.NewHandler(repo, []byte("secret_key"))

	r := gin.New()
	h.RegisterRouter(r.Group(""))

	// 1. Test Register
	regPayload := map[string]string{
		"username": "tester",
		"password": "pwd",
	}
	body, _ := json.Marshal(regPayload)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected register code 201, got %d", w.Code)
	}

	// 2. Test Login
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected login code 200, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == "" {
		t.Fatal("expected token to be returned, got empty")
	}
}
