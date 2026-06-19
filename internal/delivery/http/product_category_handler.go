package http

import (
	"net/http"
	"strconv"
	"study-golang-backend/internal/delivery/http/middleware"
	"study-golang-backend/internal/domain/entity"
	"study-golang-backend/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	repo repository.CategoryRepository
}

func NewCategoryHandler(repo repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

func (h *CategoryHandler) RegisterRouter(rg *gin.RouterGroup) {
	categories := rg.Group("/product-categories")
	{
		categories.POST("/", middleware.RequirePermission("admin"), h.CreateCategory)
		categories.GET("/", middleware.RequirePermission("admin, seller, customer"), h.GetAllCategories)
		categories.GET("/:id", middleware.RequirePermission("admin, seller, customer"), h.GetCategoryByID)
		categories.PUT("/:id", middleware.RequirePermission("admin"), h.UpdateCategory)
		categories.DELETE("/:id", middleware.RequirePermission("admin"), h.DeleteCategory)
	}
}

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Adds a new category to the inventory
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        request  body      CreateCategoryPayload  true  "Category details"
// @Success      201      {object}  Category
// @Failure      400      {object}  map[string]string "Bad Request"
// @Failure      500      {object}  map[string]string "Internal Server Error"
// @Security     BearerAuth
// @Router       /product-categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var payload entity.Category

	// Check Error from request
	// c.ShouldBindJSON is set value of payload but if error return to err not nil
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})	
		return
	}

	// Create
	newCategory, err := h.repo.Create(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": err.Error()})
		return
	}

	//Response
	c.JSON(http.StatusCreated, gin.H{
		"status":   "Created",
		"message":  "Category created successfully",
		"category": newCategory,
	})
}

// GetAllCategories godoc
// @Summary      Get all categories
// @Description  Retrieves a list of all categories in the inventory
// @Tags         categories
// @Produce      json
// @Success      200  {array}   Category
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Security     BearerAuth
// @Router       /product-categories [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "categories": categories})
}

// GetCategoryByID godoc
// @Summary      Get category by ID
// @Description  Retrieves a single category by its unique ID
// @Tags         categories
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  Category
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      404  {object}  map[string]string "Not Found"
// @Security     BearerAuth
// @Router       /product-categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid ID format"})
		return
	}
	category, err := h.repo.GetByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "category": category})
}

// UpdateCategory godoc
// @Summary      Update an category
// @Description  Updates the details of an existing category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Param        request  body      CategoryPayload  true  "Updated category details"
// @Success      200  {object}  Category
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      404  {object}  map[string]string "Not Found"
// @Security     BearerAuth
// @Router       /product-categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	// Check id is number
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid ID format"})
		return
	}

	// Check payload
	var payload entity.Category
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
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "message": "Category updated successfully", "category": updated})
}

// DeleteCategory godoc
// @Summary      Delete an category
// @Description  Removes an category from the inventory by its ID
// @Tags         categories
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {string}  string "Category deleted successfully"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      404  {object}  map[string]string "Not Found"
// @Security     BearerAuth
// @Router       /product-categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "message": "Category deleted successfully"})
}
