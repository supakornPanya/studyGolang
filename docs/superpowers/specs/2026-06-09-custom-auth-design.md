# Custom Authentication & Authorization System Design (Simplified)

This document details the simplified authentication and authorization system featuring local registration, username/password login, stateless JWT validation containing user permissions, and custom permission guards.

---

## 1. Database Schema (`User` Model)

The user system is stored in a single table. The `User` model contains credentials alongside boolean flags representing granular permissions:

```go
type User struct {
	ID           string `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"not null" json:"-"`
	CanRead      bool   `gorm:"default:true" json:"can_read"`
	CanWrite     bool   `gorm:"default:false" json:"can_write"`
	CanUpdate    bool   `gorm:"default:false" json:"can_update"`
	CanDelete    bool   `gorm:"default:false" json:"can_delete"`
}
```

---

## 2. API Endpoints

### A. `POST /register`
* **Purpose**: Register a new user with safe default permissions.
* **Payload**:
  ```json
  {
    "username": "example",
    "password": "securepassword"
  }
  ```
* **Process**:
  1. Hash the password using `bcrypt`.
  2. Save the user to the database with default permissions (`CanRead = true`, others `false`).
* **Response**: `201 Created` with a success message:
  ```json
  {
    "status": "Created",
    "message": "User registered successfully"
  }
  ```

### B. `POST /login`
* **Purpose**: Authenticate user credentials and return a signed JWT token containing their permissions.
* **Payload**:
  ```json
  {
    "username": "example",
    "password": "securepassword"
  }
  ```
* **Process**:
  1. Retrieve user from the database by `username`.
  2. Verify the submitted password against the `PasswordHash` using `bcrypt`.
  3. Generate and return a signed JWT token containing the user's identity and permission claims.
* **Response**: `200 OK` returning:
  ```json
  {
    "status": "Ok",
    "token": "<JWT Token>"
  }
  ```

---

## 3. JWT Claims

* **Session Token Claims** (includes embedded permissions):
  ```go
  type Claims struct {
  	Username  string `json:"username"`
  	UserID    string `json:"user_id"`
  	CanRead   bool   `json:"can_read"`
  	CanWrite  bool   `json:"can_write"`
  	CanUpdate bool   `json:"can_update"`
  	CanDelete bool   `json:"can_delete"`
  	jwt.RegisteredClaims
  }
  ```

---

## 4. Auth Middleware & Route Guards

1. **`AuthMiddleware`**:
   * Intercepts incoming requests on protected endpoints.
   * Extracts the JWT from `Authorization: Bearer <JWT>`.
   * Verifies the token signature and expiration.
   * Inject claims into the Gin context (e.g. `c.Set("username", claims.Username)`, `c.Set("can_read", claims.CanRead)`).

2. **Route Guards (e.g. `RequirePermission(permissionField string)`)**:
   * Wrap handlers to read the permission flags stored in the Gin context.
   * If the specified flag (e.g. `can_write`) is `false`, return `403 Forbidden`.
