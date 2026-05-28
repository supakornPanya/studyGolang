package user

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"study-golang-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo   Repository
	secret []byte
}

func NewHandler(repo Repository, secret []byte) *Handler {
	return &Handler{
		repo:   repo,
		secret: secret,
	}
}

type authRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Register(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to secure password"})
		return
	}

	b := make([]byte, 8)
	rand.Read(b)
	id := hex.EncodeToString(b)

	u := &User{
		ID:           id,
		Username:     req.Username,
		PasswordHash: hash,
	}

	if err := h.repo.Create(u); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.repo.GetByUsername(req.Username)
	if err != nil || !CheckPasswordHash(req.Password, u.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := auth.GenerateToken(u.Username, h.secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) RegisterRouter(r *gin.RouterGroup) {
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
}
