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

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user profile with a hashed password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      authRequest  true  "User registration credentials"
// @Success      201      {string}  string "User registered successfully"
// @Failure      400      {object}  map[string]string "error: bad request"
// @Failure      409      {object}  map[string]string "error: user conflict"
// @Router       /register [post]
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

// Login godoc
// @Summary      Login user
// @Description  Authenticates a user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      authRequest  true  "User login credentials"
// @Success      200      {string}  string "jwt_token_string"
// @Failure      400      {object}  map[string]string "error: bad request"
// @Failure      410      {object}  map[string]string "error: invalid username or password"
// @Router       /login [post]
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
