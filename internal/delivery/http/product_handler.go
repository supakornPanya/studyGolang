package http

import (
	"net/http"
	"strconv"
	"study-golang-backend/internal/delivery/http/middleware"
	"study-golang-backend/internal/domain/entity"
	"study-golang-backend/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// getOwnerByID get user_id from product
func (h *ProductHandler) getOwnerByID(id uint64) (string, error) {
	product, err := h.repo.GetByID(id)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(product.CreatedBy, 10), nil
}

func (h *ProductHandler) RegisterRouter(rg *gin.RouterGroup) {
	product := rg.Group("/products")
	{
		product.POST("/", middleware.RequirePermission("admin, seller"), h.CreateProduct)
		product.GET("/", middleware.RequirePermission("admin, seller, customer"), h.GetAllProducts)
		product.GET("/:id", middleware.RequirePermission("admin, seller, customer"), h.GetProductByID)
		product.PUT("/:id", middleware.RequirePermission("admin, seller"), middleware.RequireOwnership("seller", h.getOwnerByID), h.UpdateProduct)
		product.DELETE("/:id", middleware.RequirePermission("admin, seller"), middleware.RequireOwnership("seller", h.getOwnerByID), h.DeleteProduct)
	}
}

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Adds a new product to the inventory
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request  body      CreateProductPayload  true  "Product details"
// @Success      201      {object}  Product
// @Failure      400      {object}  map[string]string "Bad Request"
// @Failure      500      {object}  map[string]string "Internal Server Error"
// @Security     BearerAuth
// @Router       /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var payload entity.Product

	// Check Error from request
	// c.ShouldBindJSON is set value of payload but if error return to err not nil
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	// Get userID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "User ID not found in context"})
		return
	}
	userIDStrVal, _ := userIDStr.(string)
	payload.CreatedBy, _ = strconv.ParseUint(userIDStrVal, 10, 64)

	// Create
	newProduct, err := h.repo.Create(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": err.Error()})
		return
	}

	//Response
	c.JSON(http.StatusCreated, gin.H{
		"status":   "Created",
		"message":  "Product created successfully",
		"products": newProduct,
	})
}

// GetAllProducts godoc
// @Summary      Get all products
// @Description  Retrieves a list of all products in the inventory
// @Tags         products
// @Produce      json
// @Success      200  {array}   Product
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Security     BearerAuth
// @Router       /products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "products": products})
}

// GetProductByID godoc
// @Summary      Get product by ID
// @Description  Retrieves a single product by its unique ID
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  Product
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      404  {object}  map[string]string "Not Found"
// @Security     BearerAuth
// @Router       /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid ID format"})
		return
	}
	product, err := h.repo.GetByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "products": product})
}

// UpdateProduct godoc
// @Summary      Update an product
// @Description  Updates the details of an existing product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Param        request  body      ProductPayload  true  "Updated product details"
// @Success      200  {object}  Product
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      404  {object}  map[string]string "Not Found"
// @Security     BearerAuth
// @Router       /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Check id is number
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid ID format"})
		return
	}

	// Check payload
	var payload entity.Product
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": err.Error()})
		return
	}

	// Update
	updated, err := h.repo.Update(id, payload)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "message": "Product updated successfully", "products": updated})
}

// DeleteProduct godoc
// @Summary      Delete an product
// @Description  Removes an product from the inventory by its ID
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {string}  string "Product deleted successfully"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      404  {object}  map[string]string "Not Found"
// @Security     BearerAuth
// @Router       /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid ID format"})
		return
	}
	if err := h.repo.Delete(uint64(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "message": "Product deleted successfully"})
}
