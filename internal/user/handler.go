package user

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"study-golang-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

// Handler dependency
type Handler struct {
	repo Repository
	secret []byte
}

// Dependency Injection
func NewHandler(repo Repository, secret []byte) *Handler {
	return &Handler{
		repo: repo,
		secret: secret,
	}
}

// RegisterRouter
func (h *Handler) RegisterRouter(r *gin.RouterGroup) {
	r.POST("register", h.Register)
	r.POST("login", h.Login)
}

// Request body payload for Register/Login
type authRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register handler
// @Summary      Register a new user
// @Description  Creates a new user with username and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      authRequest  true  "User credentials"
// @Success      201      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /register [post]
func (h *Handler) Register(c *gin.Context) {
	var req authRequest

	// Bind and validate request body
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	// Get Hash password
	hash, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": err.Error()})
		return
	}
	
	// Generate random unique ID
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	id := base64.StdEncoding.EncodeToString(b)

	// Create User entity
	user := &User{
		ID: id,
		Username: req.Username,
		PasswordHash: hash,
		CanRead: true,
		CanWrite: false,
		CanUpdate: false,
		CanDelete: false,
	}

	// Save user
	if err := h.repo.Create(user); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"status": "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "Created",
		"message": "User created successfully",
	})
}

// Login handler 
// @Summary      Login
// @Description  Logs in a user and returns an authentication token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      authRequest  true  "User credentials"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /login [post]
func (h *Handler) Login(c *gin.Context) {
	var req authRequest

	// Validate Json
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
			"error": err.Error(),
		})
		return
	}

	// Fetch user from DB -> Check password
	u, err := h.repo.GetByUsername(req.Username)
	if err != nil || !CheckPassword(req.Password, u.PasswordHash){
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "Unauthorized",
			"message": "Invalid username or password",
		})
		return
	}

	// Generate token from json
	token, err := auth.GenerateToken(u.Username, u.ID, u.CanRead, u.CanWrite, u.CanUpdate, u.CanDelete, h.secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Ok",
		"message": "Login successfully",
		"data": gin.H{
			"token": token,
		},
	})
}