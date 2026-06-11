# Spec: JWT User Authentication

This document outlines the design for implementing JWT User Authentication in the Go/Gin backend project.

## Goal
Add user registration, authentication, and JWT verification middleware to secure the existing item endpoints.

## Proposed Changes

### 1. `internal/user` Package

#### [NEW] [model.go](file:///c:/studyGolang/internal/user/model.go)
Contains the user data models and repository interface:
- `User` struct (ID, Username, PasswordHash)
- `Repository` interface:
  - `Create(user *User) error`
  - `GetByUsername(username string) (*User, error)`
- `NewMemoryRepository()` returning an in-memory repository implementation for users.

#### [NEW] [handler.go](file:///c:/studyGolang/internal/user/handler.go)
Contains HTTP handlers for user operations:
- `Register(c *gin.Context)`: Creates a new user, hashing their password with bcrypt.
- `Login(c *gin.Context)`: Validates credentials and generates a JWT.
- `RegisterRouter(r *gin.RouterGroup)`: Registers the `/register` and `/login` endpoints.

### 2. `internal/auth` Package

#### [NEW] [middleware.go](file:///c:/studyGolang/internal/auth/middleware.go)
Contains JWT authentication middleware:
- `AuthMiddleware(secret []byte) gin.HandlerFunc`: Inspects the `Authorization` header, validates the JWT, and sets the authenticated username in the Gin context.

#### [NEW] [token.go](file:///c:/studyGolang/internal/auth/token.go)
Contains functions for token lifecycle management:
- `GenerateToken(username string, secret []byte) (string, error)`
- `ValidateToken(tokenString string, secret []byte) (string, error)`

### 3. Main Application Entrypoint

#### [MODIFY] [main.go](file:///c:/studyGolang/cmd/api/main.go)
- Initialize the user repository and handler.
- Register user authentication routes (`POST /register` and `POST /login`).
- Apply the `AuthMiddleware` to the item routes to secure them.

---

## Verification Plan

### Automated Tests
We will implement unit tests using Go's standard `testing` package and Gin's test utilities (`net/http/httptest`):
- Test password hashing and verification.
- Test JWT token generation and validation.
- Test `AuthMiddleware` denies unauthorized requests and passes authorized ones.
- Test registration & login HTTP handler endpoints.

Run tests using:
```bash
go test ./...
```
