package main

import (
	"log"
	"study-golang-backend/internal/auth"
	"study-golang-backend/internal/cart"
	"study-golang-backend/internal/db"
	"study-golang-backend/internal/user"

	_ "study-golang-backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Study Golang Backend API
// @version         1.0
// @description     This is a sample backend API using Gin, PostgreSQL, and Redis.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Type "Bearer " followed by your JWT token.
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
	err = database.AutoMigrate(&cart.Item{})
	if err != nil {
		log.Fatal("Failed to auto migrate database", err)
	}

	// Secret Key For JWT
	jwtSecret := []byte("yeahhhh_this_is_super_super_ummmm_secret_key_right?")

	//Init Gin router
	r := gin.Default()

	// Swagger Documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initial User
	userRepo := user.NewPostgresRepository(database)
	userHandler := user.NewHandler(userRepo, jwtSecret)
	userHandler.RegisterRouter(r.Group(""))

	// Initial Cart
	rdbClient := db.InitRedis()
	postresCartRepo := cart.NewPostgresRepository(database)
	cartRepo := cart.NewCachedRepository(postresCartRepo, rdbClient)
	// Initial Handler & inject initial Data(cartRepo)
	cartHandler := cart.NewHandler(cartRepo)

	// Register Router secured by AuthMiddleware
	protected := r.Group("")
	protected.Use(auth.AuthMiddleware(jwtSecret))
	cartHandler.RegisterRouter(protected)

	r.Run(":8080")
}
