package cart

import (
	"net/http"
	"strconv"
	"study-golang-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRouter(rg *gin.RouterGroup) {
	items := rg.Group("/items")
	{
		items.POST("/", auth.RequirePermission("can_write"), h.CreateItem)
		items.GET("/", auth.RequirePermission("can_read"), h.GetAllItems)
		items.GET("/:id", auth.RequirePermission("can_read"), h.GetItemByID)
		items.PUT("/:id", auth.RequirePermission("can_update"), h.UpdateItem)
		items.DELETE("/:id", auth.RequirePermission("can_delete"), h.DeleteItem)
	}
}

// CreateItem godoc
// @Summary      Create a new item
// @Description  Adds a new item to the inventory
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        request  body      CreateItemPayload  true  "Item details"
// @Success      201      {object}  Item
// @Failure      400      {object}  map[string]string "Bad Request"
// @Failure      500      {object}  map[string]string "Internal Server Error"
// @Security     BearerAuth
// @Router       /items [post]
func (h *Handler) CreateItem(c *gin.Context) {
	var payload CreateItemPayload

	// Check Error from request
	// c.ShouldBindJSON is set value of payload but if error return to err not nil
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	newItem, err := h.repo.Create(payload.SKU, payload.Name, payload.Quantity, payload.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": err.Error()})
		return
	}

	//Response
	c.JSON(http.StatusCreated, gin.H{
		"status":  "Created",
		"message": "Item created successfully",
		"data":    newItem,
	})
}

// GetAllItems godoc
// @Summary      Get all items
// @Description  Retrieves a list of all items in the inventory
// @Tags         items
// @Produce      json
// @Success      200  {array}   Item
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Security     BearerAuth
// @Router       /items [get]
func (h *Handler) GetAllItems(c *gin.Context) {
	items, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "items": items})
}

// GetItemByID godoc
// @Summary      Get item by ID
// @Description  Retrieves a single item by its unique ID
// @Tags         items
// @Produce      json
// @Param        id   path      int  true  "Item ID"
// @Success      200  {object}  Item
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      404  {object}  map[string]string "Not Found"
// @Security     BearerAuth
// @Router       /items/{id} [get]
func (h *Handler) GetItemByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid ID format"})
		return
	}
	item, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "item": item})
}

// UpdateItem godoc
// @Summary      Update an item
// @Description  Updates the details of an existing item
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Item ID"
// @Param        request  body      UpdateItemPayload  true  "Updated item details"
// @Success      200  {object}  Item
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      404  {object}  map[string]string "Not Found"
// @Security     BearerAuth
// @Router       /items/{id} [put]
func (h *Handler) UpdateItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid ID format"})
		return
	}
	var payload UpdateItemPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": err.Error()})
		return
	}
	updated, err := h.repo.Update(id, payload.SKU, payload.Name, payload.Quantity, payload.Price)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "message": "Item updated successfully", "data": updated})
}

// DeleteItem godoc
// @Summary      Delete an item
// @Description  Removes an item from the inventory by its ID
// @Tags         items
// @Produce      json
// @Param        id   path      int  true  "Item ID"
// @Success      200  {string}  string "Item deleted successfully"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      404  {object}  map[string]string "Not Found"
// @Security     BearerAuth
// @Router       /items/{id} [delete]
func (h *Handler) DeleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid ID format"})
		return
	}
	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "message": "Item deleted successfully"})
}
