# Custom Authentication & Authorization System Design

This document details the direct authentication system featuring local registration, a two-tier token exchange (Parent Token and Local Session tokens), stateless JWT validation, and custom permission guards.

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
* **Purpose**: Register a new user and return a signed Parent Token (simulating the external gateway).
* **Payload**:
  ```json
  {
    "username": "example",
    "password": "securepassword"
  }
  ```
* **Process**:
  1. Hashing the password using `bcrypt`.
  2. Saving the user with default permissions (`CanRead = true`, others `false`).
  3. Generating a signed **Parent Token** (JWT) containing `username` and `id` claims.
* **Response**: `201 Created` with the parent token.

### B. `POST /token`
* **Purpose**: Exchange the parent token for local session tokens.
* **Headers**: `X-Access-Token: <Parent JWT>`
* **Process**:
  1. Structurally decode and verify the Parent Token signature.
  2. Fetch the user details and permissions from the database.
  3. Generate and return a local **Access Token** and **Refresh Token** pair.
* **Response**: `200 OK` returning:
  ```json
  {
    "access_token": "<Local Access JWT>",
    "refresh_token": "<Local Refresh JWT>"
  }
  ```

### C. `POST /refresh-token`
* **Purpose**: Rotate expired access tokens using the local refresh token.
* **Payload**:
  ```json
  {
    "refresh_token": "<Local Refresh JWT>"
  }
  ```
* **Process**:
  1. Verify the signature, expiry, and check if `token_use` is `"refresh"`.
  2. If expired or invalid, return `401 Unauthorized`.
  3. If valid, issue a new access token and a rotated refresh token.
* **Response**: `200 OK` with the new token pair.

---

## 3. JWT Claims & Tokens

* **Parent Token Claims**:
  ```go
  type ParentClaims struct {
  	Username string `json:"username"`
  	jwt.RegisteredClaims
  }
  ```
* **Local Session Access Token Claims** (includes embedded permissions):
  ```go
  type AccessClaims struct {
  	Username   string `json:"username"`
  	UserID     string `json:"user_id"`
  	CanRead    bool   `json:"can_read"`
  	CanWrite   bool   `json:"can_write"`
  	CanUpdate  bool   `json:"can_update"`
  	CanDelete  bool   `json:"can_delete"`
  	TokenUse   string `json:"token_use"` // Must be "access"
  	jwt.RegisteredClaims
  }
  ```
* **Local Session Refresh Token Claims**:
  ```go
  type RefreshClaims struct {
  	Username string `json:"username"`
  	TokenUse string `json:"token_use"` // Must be "refresh"
  	jwt.RegisteredClaims
  }
  ```

---

## 4. Auth Middleware & Route Guards

1. **`AuthMiddleware`**:
   * Intercepts incoming requests on protected endpoints.
   * Extracts the local Access Token from `Authorization: Bearer <Access JWT>`.
   * Verifies the token signature, check expiration, and validates `TokenUse == "access"`.
   * Inject claims into the Gin context (e.g. `c.Set("claims", claims)`).

2. **Route Guards (e.g. `RequirePermission(permissionName string)`)**:
   * Wrap handlers to read the permission flags stored in the Gin context.
   * If the flag (e.g. `CanWrite`) is `false`, immediately return `403 Forbidden`.
