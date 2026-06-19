package http

import (
	"net/http"
	"strconv"
	"study-golang-backend/internal/delivery/http/middleware"
	"study-golang-backend/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	repo repository.CartRepository
}

func NewCartHandler(repo repository.CartRepository) *CartHandler {
	return &CartHandler{repo: repo}
}

func (h *CartHandler) RegisterRouter(rg *gin.RouterGroup) {
	cart := rg.Group("/cart")
	{
		cart.POST("/add", middleware.RequirePermission("admin, seller, customer"), h.AddCart)
		cart.GET("/", middleware.RequirePermission("admin, seller, customer"), h.GetCart)
	}
}

// AddCart godoc
// @Summary      Add a new product to the cart
// @Description  Adds a new product to the cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        request  body      entity.CartItems  true  "Product details"
// @Success      201      {object}  entity.CartItems
// @Failure      400      {object}  map[string]string "Bad Request"
// @Failure      500      {object}  map[string]string "Internal Server Error"
// @Security     BearerAuth
// @Router       /cart/add [post]
func (h *CartHandler) AddCart(c *gin.Context) {
	// Get userID from context
	userIDVal, _ := c.Get("user_id")
	userIDStr := userIDVal.(string)
	userID, _ := strconv.Atoi(userIDStr)

	// Check payload
	var payload struct {
		ProductID []uint64 `json:"product_id"binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": err.Error()})
		return
	}

	// 1. Fetch current cart
	cart, err := h.repo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": err.Error()})
		return
	}

	// 2. Append to the product new to old list product
	cart.ListProductID = append(cart.ListProductID, payload.ProductID...)

	// 3. Save back to Redis
	updatedCart, err := h.repo.Save(userID, cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Ok", "message": "Product added to cart successfully", "cart": updatedCart})
}

// GetCart godoc
// @Summary      Get cart by ID
// @Description  Retrieves a single cart by its unique ID
// @Tags         cart
// @Produce      json
// @Param        id   path      int  true  "Cart ID"
// @Success      200  {object}  entity.CartItems
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      404  {object}  map[string]string "Not Found"
// @Security     BearerAuth
// @Router       /cart/{id} [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	// Get userID from context
	userIDVal, _ := c.Get("user_id")
	userIDStr := userIDVal.(string)
	userID, _ := strconv.Atoi(userIDStr)

	// Get Cart by userID
	cart, err := h.repo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "cart": cart})
}
