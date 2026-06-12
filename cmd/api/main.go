package main

import (
	"log/slog"
	"os"
	"study-golang-backend/internal/auth"
	"study-golang-backend/internal/db"
	"study-golang-backend/internal/logger"
	"study-golang-backend/internal/domain/entity"
	"study-golang-backend/internal/domain/repository"
	"study-golang-backend/internal/delivery/http"

	_ "study-golang-backend/docs"

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
	logger.InitLogger()
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found")
	}

	// Initialize Database
	// Auto Migrate database: Check has this table yes->skip, no->create
	database := db.InitDB()
	err = database.AutoMigrate(&entity.User{})
	if err != nil {
		slog.Error("Failed to auto migrate database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	err = database.AutoMigrate(&entity.Item{})
	if err != nil {
		slog.Error("Failed to auto migrate database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Secret Key For JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		slog.Error("JWT_SECRET not found in .env file")
		os.Exit(1)
	}
	jwtSecret := []byte(secret)

	//Init Gin router & Register logger middleware
	r := gin.New()
	r.Use(logger.LoggerMiddleware(), gin.Recovery())

	// Swagger Documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initial User
	userRepo := repository.NewUserPostgreRepository(database)
	userHandler := http.NewUserHandler(userRepo, jwtSecret)
	userHandler.RegisterRouter(r.Group(""))

	// Initial Cart
	rdbClient := db.InitRedis()
	postresCartRepo := repository.NewCartPostgresRepository(database)
	cartRepo := repository.NewCartCachedRepository(postresCartRepo, rdbClient)
	cartHandler := http.NewCartHandler(cartRepo)

	// Register Router secured by AuthMiddleware
	protected := r.Group("")
	protected.Use(auth.AuthMiddleware(jwtSecret))
	cartHandler.RegisterRouter(protected)

	r.Run(":8080")
}
