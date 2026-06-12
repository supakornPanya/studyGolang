# Structured Logging Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Complete the migration to Go's structured `log/slog` logging library by replacing the remaining standard `"log"` library usages in `internal/db/db.go` and `cmd/api/main.go`.

**Architecture:** Initialize Go's `log/slog` at start, register custom Gin logging middleware, and use `slog` globally for system logs.

**Tech Stack:** Go (slog), Gin

---

### Task 1: Update db.go to use slog
Modify `internal/db/db.go` to import `"log/slog"` instead of `"log"` and use structured log methods for error reporting and success notifications.

**Files:**
- Modify: `c:/studyGolang/internal/db/db.go`

- [ ] **Step 1: Replace log references with slog**
  Replace imports and logging code in `internal/db/db.go`:
  ```diff
  package db

  import (
  	"fmt"
- 	"log"
  	"os"
+ 	"log/slog"

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
- 		log.Fatalf("Failed to connect to database: %v", err)
+ 		slog.Error("Failed to connect to database", slog.Any("error", err))
+ 		os.Exit(1)
  	}

- 	log.Println("Successfully connected to database")
+ 	slog.Info("Successfully connected to database")

  	return db
  }
  ```

- [ ] **Step 2: Verify compilation**
  Run: `go build ./internal/db`
  Expected: Command finishes successfully with exit status 0 (no errors).

---

### Task 2: Update main.go to use slog
Modify `cmd/api/main.go` to import `"log/slog"` instead of `"log"` and use structured log methods for startup errors.

**Files:**
- Modify: `c:/studyGolang/cmd/api/main.go`

- [ ] **Step 1: Replace log references with slog**
  Replace standard log calls in `cmd/api/main.go`:
  ```diff
  package main

  import (
- 	"log"
  	"os"
+ 	"log/slog"
  	"study-golang-backend/internal/auth"
  	"study-golang-backend/internal/cart"
  	"study-golang-backend/internal/db"
  	"study-golang-backend/internal/logger"
  	"study-golang-backend/internal/user"
  ...
  func main() {
  	logger.InitLogger()
  	err := godotenv.Load()
  	if err != nil {
- 		log.Println("Warning: No .env file found")
+ 		slog.Warn("No .env file found")
  	}

  	// Initialize Database
  	database := db.InitDB()
  	// Auto Migrate database: Check has this table yes->skip, no->create
  	err = database.AutoMigrate(&user.User{})
  	if err != nil {
- 		log.Fatal("Failed to auto migrate database", err)
+ 		slog.Error("Failed to auto migrate database", slog.Any("error", err))
+ 		os.Exit(1)
  	}
  	err = database.AutoMigrate(&cart.Item{})
  	if err != nil {
- 		log.Fatal("Failed to auto migrate database", err)
+ 		slog.Error("Failed to auto migrate database", slog.Any("error", err))
+ 		os.Exit(1)
  	}

  	// Secret Key For JWT
  	secret := os.Getenv("JWT_SECRET")
  	if secret == "" {
- 		log.Fatal("JWT_SECRET not found in .env file")
+ 		slog.Error("JWT_SECRET not found in .env file")
+ 		os.Exit(1)
  	}
  ```

- [ ] **Step 2: Verify compilation**
  Run: `go build ./cmd/api`
  Expected: Command finishes successfully with exit status 0 (no errors).

---

### Task 3: Run and Verify Logs
Run the application and verify that all startup and HTTP logs are output using the structured slog text format.

**Files:**
- Test: none (run backend)

- [ ] **Step 1: Start the backend**
  Run: `go run ./cmd/api/main.go`
  Expected: Prints structured slog text output for database connection, e.g.:
  `time=... level=INFO msg="Successfully connected to database"`

- [ ] **Step 2: Send a request to trigger middleware logging**
  Run (in another terminal or via curl/Invoke-RestMethod):
  `Invoke-RestMethod -Uri "http://localhost:8080/items/" -Method GET`
  Expected: Logs a request message structured like:
  `time=... level=INFO msg="gin request" method=GET path=/items/ query="" status=200 latency=... client_ip=::1 user_agent=...`
