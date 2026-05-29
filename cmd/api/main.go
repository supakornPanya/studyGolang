package main

import (
	"log"
	"net/http"
	"study-golang-backend/internal/auth"
	"study-golang-backend/internal/db"
	"study-golang-backend/internal/item"
	"study-golang-backend/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	// Initialize Database
	database := db.InitDB()
	// Auto Migrate database: Check has this table yes->skip, no->create
	err = database.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatal("Failed to auto migrate database", err)
	}
	err = database.AutoMigrate(&item.Item{})
	if err != nil {
		log.Fatal("Failed to auto migrate database", err)
	}

	// JWT Secret Key: Must be stored securely in Production
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
	userRepo := user.NewPostgresRepository(database)
	userHandler := user.NewHandler(userRepo, jwtSecret)
	userHandler.RegisterRouter(r.Group(""))

	// Initial Item
	rdbClient := db.InitRedis()
	postresItemRepo := item.NewPostgresRepository(database)
	itemRepo := item.NewCachedRepository(postresItemRepo, rdbClient)
	// Initial Handler & inject initial Data(itemRepo)
	itemHandler := item.NewHandler(itemRepo)

	// Register Router secured by AuthMiddleware
	protected := r.Group("")
	protected.Use(auth.AuthMiddleware(jwtSecret))
	itemHandler.RegisterRouter(protected)

	r.Run(":8080")
}
