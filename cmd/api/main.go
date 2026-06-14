package main

import (
	"log/slog"
	"os"

	"study-golang-backend/internal/db"
	delivery "study-golang-backend/internal/delivery/http"
	"study-golang-backend/internal/domain/entity"
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
	fx.New(
		// 1. Group all our modules together
		db.Module,
		repository.Module,
		delivery.Module,
		// 2. Invoke setup tasks (Logger, Env, Migrations, Swagger)
		fx.Invoke(
			setupEnvAndLogger,
			runMigrations,
			setupSwagger,
		),
	).Run()
}
// setupEnvAndLogger loads environment variables and sets up the logger
func setupEnvAndLogger() {
	logger.InitLogger()
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found")
	}
}
// runMigrations automatically migrates the GORM tables
func runMigrations(database *gorm.DB) {
	err := database.AutoMigrate(&entity.User{})
	if err != nil {
		slog.Error("Failed to auto migrate User database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	err = database.AutoMigrate(&entity.Item{})
	if err != nil {
		slog.Error("Failed to auto migrate Item database", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
// setupSwagger registers the Swagger endpoint
func setupSwagger(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}