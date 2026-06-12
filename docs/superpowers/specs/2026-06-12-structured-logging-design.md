# Structured Logging Design Spec

Design for transitioning the study-golang-backend from basic standard library logging to Go's structured `log/slog` framework.

## Proposed Design

We introduce a structured logging system using the native Go standard library package `log/slog`. This enables consistent logging formats (text in local development, JSON in production) and formats HTTP requests properly.

### 1. Logger Setup (`internal/logger/logger.go`)
Initialize a global default logger configurable by `APP_ENV`:
- **Development/Local:** Colored or plain structured Text Handler (`slog.NewTextHandler`) at `Debug` level.
- **Production:** JSON Handler (`slog.NewJSONHandler`) at `Info` level.

### 2. Custom Gin Middleware (`internal/logger/middlerware.go`)
Intercept incoming requests to log:
- `method`, `path`, `query`, `status`, `latency`, `client_ip`, `user_agent`.
- Choose severity level dynamically based on response status code:
  - Status `>= 500` -> `ERROR`
  - Status `>= 400` -> `WARN`
  - Status `< 400` -> `INFO`

### 3. Service Updates
We migrate all remaining `log.Print` / `log.Fatal` / `log.Println` occurrences to `slog` methods:
- `cmd/api/main.go`
- `internal/db/db.go`

---

## Verification Plan

### Manual Verification
1. Run application: `go run .\cmd\api\main.go`
2. Perform mock HTTP requests to verify structured request logging format:
   - Success requests (`200 OK`) should log as `LevelInfo`.
   - Client errors (`403 Forbidden`, `400 Bad Request`) should log as `LevelWarn`.
   - Verify connection messages from PostgreSQL/Redis connection.
