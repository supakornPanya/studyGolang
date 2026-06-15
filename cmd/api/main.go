package main

import (
	"log/slog"
	"os"

	"study-golang-backend/internal/db"
	delivery "study-golang-backend/internal/delivery/http"
	"study-golang-backend/internal/infrastructure/logger"
	"study-golang-backend/internal/repository"

	_ "study-golang-backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"go.uber.org/fx"
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
	// 1. load environment variables
	logger.InitLogger()
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found")
	}

	fx.New(
		// 2. Group all our modules together
		db.Module,
		repository.Module,
		delivery.Module,
		// 3. Invoke setup tasks (Migrations, Swagger)
		fx.Invoke(
			runMigrations,
			setupSwagger,
		),
	).Run()
}

// runMigrations automatically migrates the GORM tables
func runMigrations(database *gorm.DB) {
	if err := db.RunMigrations(database); err != nil {
		slog.Error("Failed to run database migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

// setupSwagger registers the Swagger endpoint
func setupSwagger(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}