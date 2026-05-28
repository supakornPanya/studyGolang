package main

import (
	"net/http"
	"study-golang-backend/internal/auth"
	"study-golang-backend/internal/item"
	"study-golang-backend/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	jwtSecret := []byte("my_super_secret_dev_key")

	//Init Gin router
	r := gin.Default()

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Ok",
			"message": "Go Backend with Gin is running!",
		})
	})

	// Setup User authentication
	userRepo := user.NewMemoryRepository()
	userHandler := user.NewHandler(userRepo, jwtSecret)
	userHandler.RegisterRouter(r.Group(""))

	// Initial example Data
	itemRepo := item.NewMemoryRepository()
	// Initial Handler & inject initial Data(itemRepo)
	itemHandler := item.NewHandler(itemRepo)

	// Register Router secured by AuthMiddleware
	protected := r.Group("")
	protected.Use(auth.AuthMiddleware(jwtSecret))
	itemHandler.RegisterRouter(protected)

	r.Run(":8080")
}

