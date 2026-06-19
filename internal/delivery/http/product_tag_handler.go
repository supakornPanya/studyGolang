package http

import (
	"net/http"
	"strconv"
	"study-golang-backend/internal/delivery/http/middleware"
	"study-golang-backend/internal/domain/entity"
	"study-golang-backend/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	repo repository.TagRepository
}

func NewTagHandler(repo repository.TagRepository) *TagHandler {
	return &TagHandler{repo: repo}
}

func (h *TagHandler) RegisterRouter(rg *gin.RouterGroup) {
	tags := rg.Group("/product-tags")
	{
		tags.POST("/", middleware.RequirePermission("admin"), h.CreateTag)
		tags.GET("/", middleware.RequirePermission("admin, seller, customer"), h.GetAllTags)
		tags.GET("/:id", middleware.RequirePermission("admin, seller, customer"), h.GetTagByID)
		tags.PUT("/:id", middleware.RequirePermission("admin"), h.UpdateTag)
		tags.DELETE("/:id", middleware.RequirePermission("admin"), h.DeleteTag)
	}
}

// CreateTag godoc
// @Summary      Create a new tag
// @Description  Adds a new tag to the inventory
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        request  body      CreateTagPayload  true  "Tag details"
// @Success      201      {object}  Tag
// @Failure      400      {object}  map[string]string "Bad Request"
// @Failure      500      {object}  map[string]string "Internal Server Error"
// @Security     BearerAuth
// @Router       /product-tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var payload entity.Tag

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
	newTag, err := h.repo.Create(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": err.Error()})
		return
	}

	//Response
	c.JSON(http.StatusCreated, gin.H{
		"status":   "Created",
		"message":  "Tag created successfully",
		"tag":      newTag,
	})
}

// GetAllTags godoc
// @Summary      Get all tags
// @Description  Retrieves a list of all tags in the inventory
// @Tags         tags
// @Produce      json
// @Success      200  {array}   Tag
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Security     BearerAuth
// @Router       /product-tags [get]
func (h *TagHandler) GetAllTags(c *gin.Context) {
	tags, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "tags": tags})
}

// GetTagByID godoc
// @Summary      Get tag by ID
// @Description  Retrieves a single tag by its unique ID
// @Tags         tags
// @Produce      json
// @Param        id   path      int  true  "Tag ID"
// @Success      200  {object}  Tag
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      404  {object}  map[string]string "Not Found"
// @Security     BearerAuth
// @Router       /product-tags/{id} [get]
func (h *TagHandler) GetTagByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid ID format"})
		return
	}
	tag, err := h.repo.GetByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "tag": tag})
}

// UpdateTag godoc
// @Summary      Update an tag
// @Description  Updates the details of an existing tag
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Tag ID"
// @Param        request  body      TagPayload  true  "Updated tag details"
// @Success      200  {object}  Tag
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      404  {object}  map[string]string "Not Found"
// @Security     BearerAuth
// @Router       /product-tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	// Check id is number
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid ID format"})
		return
	}

	// Check payload
	var payload entity.Tag
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
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "message": "Tag updated successfully", "tag": updated})
}

// DeleteTag godoc
// @Summary      Delete an tag
// @Description  Removes an tag from the inventory by its ID
// @Tags         tags
// @Produce      json
// @Param        id   path      int  true  "Tag ID"
// @Success      200  {string}  string "Tag deleted successfully"
// @Failure      400  {object}  map[string]string "Bad Request"
// @Failure      404  {object}  map[string]string "Not Found"
// @Security     BearerAuth
// @Router       /product-tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"status": "Ok", "message": "Tag deleted successfully"})
}
