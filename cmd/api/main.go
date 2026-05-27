package main

import (
	"net/http"
	"study-golang-backend/internal/item"

	"github.com/gin-gonic/gin"
)

func main() {
	//Init Gin router
	r := gin.Default()

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Ok",
			"message": "Go Backend with Gin is running!",
		})
	})

	// Initial example Data
	itemRepo := item.NewMemoryRepository()
	// Initial Handler & inject initial Data(itemRepo)
	itemHandler := item.NewHandler(itemRepo)
	// Register Router
	itemHandler.RegisterRouter(r.Group(""))

	r.Run(":8080")
}
