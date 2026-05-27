package item

import (
	"net/http"
	"strconv"

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
		items.POST("/", h.CreateItem)
		items.GET("/", h.GetAllItems)
		items.GET("/:id", h.GetItemByID)
		items.PUT("/:id", h.UpdateItem)
		items.DELETE("/:id", h.DeleteItem)
	}
}

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

	newItem := h.repo.Create(payload.SKU, payload.Name, payload.Quantity, payload.Price)

	//Response
	c.JSON(http.StatusCreated, gin.H{
		"status":  "Created",
		"message": "Item created successfully",
		"data":    newItem,
	})
}

func (h *Handler) GetAllItems(c *gin.Context) {
	items := h.repo.GetAll()
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "items": items})
}

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
