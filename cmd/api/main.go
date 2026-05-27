package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Create a struct for items
type Item struct {
	ID       int     `json:"id"`
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// Initial item in items
var (
	items = []Item{
		{
			ID:       1,
			SKU:      "SKU001",
			Name:     "Item 1",
			Quantity: 10,
			Price:    100,
		},
		{
			ID:       2,
			SKU:      "SKU002",
			Name:     "Item 2",
			Quantity: 20,
			Price:    200,
		},
	}
	nextID = 3
)

func StatusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "Ok",
		"message": "Go Backend with Gin is running!",
	})
}

func CreateItemHandler(c *gin.Context) {
	// Create template for payload that is json can't using in code like a struct
	type CreateItemPayload struct {
		SKU      string  `json:"sku" binding:"required"`
		Name     string  `json:"name" binding:"required"`
		Quantity int     `json:"quantity" binding:"required,gt=0"`
		Price    float64 `json:"price" binding:"required,gt=0"`
	}

	var payload CreateItemPayload

	// Check Error from request
	// c.ShouldBindJSON is set value of payload but if error return to err not nil
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid request payload",
		})
		return
	}

	newItem := Item{
		ID:       nextID,
		SKU:      payload.SKU,
		Name:     payload.Name,
		Quantity: payload.Quantity,
		Price:    payload.Price,
	}
	nextID++

	//append item to items
	items = append(items, newItem)

	//Response
	c.JSON(http.StatusCreated, gin.H{
		"status":  "Created",
		"message": "Item created successfully",
		"data":    newItem,
	})
}

func GetItemsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "Ok",
		"items":  items,
	})
}

func GetItemHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	//Check error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid ID format",
		})
		return
	}

	for _, item := range items {
		if item.ID == id {
			c.JSON(http.StatusOK, gin.H{
				"status": "Ok",
				"item":   item,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  "Not Found",
		"message": "Item not found",
	})
}

func UpdateItemHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	//Check error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid ID format",
		})
		return
	}

	//Create template for payload that is json can't using in code like a struct
	type UpdateItemPayload struct {
		SKU      string  `json:"sku"`
		Name     string  `json:"name"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}

	var payload UpdateItemPayload

	//Check error from request
	// c.ShouldBindJSON is set value of payload but if error return to err not nil
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid request payload",
		})
		return
	}

	//Check exist or not
	for i, item := range items {
		if id == item.ID {
			if payload.Name != "" {
				items[i].Name = payload.Name
			}
			if payload.SKU != "" {
				items[i].SKU = payload.SKU
			}
			if payload.Quantity > 0 {
				items[i].Quantity = payload.Quantity
			}
			if payload.Price > 0 {
				items[i].Price = payload.Price
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  "Ok",
				"message": "Item updated successfully",
				"data":    items[i],
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  "Not Found",
		"message": "Item not found",
	})
}

func DeleteItemHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	//Check error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid ID format",
		})
		return
	}

	for i, item := range items {
		if item.ID == id {
			items = append(items[:i], items[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"status":  "Ok",
				"message": "Item deleted successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  "Not Found",
		"message": "Item not found",
	})
}

func main() {
	//Init Gin router
	r := gin.Default()

	// REST API
	r.GET("/status", StatusHandler)
	r.POST("/items", CreateItemHandler)
	r.GET("/items", GetItemsHandler)
	r.GET("/items/:id", GetItemHandler)
	r.PUT("/items/:id", UpdateItemHandler)
	r.DELETE("/items/:id", DeleteItemHandler)

	r.Run(":8080")
}
