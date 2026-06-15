package db

import (
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Initial Database PostgreSQL by GORM
func InitDB() *gorm.DB {
	// Get database configuration from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_MODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, dbname, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("Successfully connected to database")

	return db
}

//Run goose for migration database
//go:embed migrations/*.sql
var embedMigrations embed.FS //get all file .sql in this folder 
func RunMigrations(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("Failed to get database connection: %w", err)
	}

	goose.SetBaseFS(embedMigrations) //give all file .sql
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	slog.Info("Running database migrations...")

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return fmt.Errorf("failed to run goose up: %w", err)
	}

	slog.Info("Database migrations completed successfully")

	return nil
}