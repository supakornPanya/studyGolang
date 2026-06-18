# Product Ownership Verification Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Enforce that sellers can only create, update, and delete products that they own (which they created), and fix the broken `RequirePermission` middleware to correctly validate user roles from the JWT claims.

**Architecture:** 
1. Modify `ProductRepository.Create` to accept and record the creator's user ID (`CreatedBy`).
2. Fix `RequirePermission` middleware to extract the `role` context key and match it against required permissions.
3. Update HTTP product handlers to extract `user_id` and `role` from the context, and verify that sellers only modify/delete their own products.

**Tech Stack:** Go, Gin, GORM, PostgreSQL

---

## Proposed Changes

### Component 1: Domain Repository

#### [MODIFY] [product.go](file:///c:/studyGolang/internal/domain/repository/product.go)
Add `createdBy uint64` to `Create` signature.

### Component 2: Infrastructure Repository

#### [MODIFY] [product.go](file:///c:/studyGolang/internal/repository/product.go)
Implement updated signature of `Create` to accept `createdBy uint64` and save it.

### Component 3: Authentication Middleware

#### [MODIFY] [auth.go](file:///c:/studyGolang/internal/delivery/http/middleware/auth.go)
Correct `RequirePermission` to check role string dynamically.

### Component 4: HTTP Handlers

#### [MODIFY] [product_handler.go](file:///c:/studyGolang/internal/delivery/http/product_handler.go)
- Pass current user ID to `CreateProduct`.
- Enforce creator checks in `UpdateProduct` and `DeleteProduct`.

---

## Tasks

### Task 1: Update Product Repository Signatures

**Files:**
- Modify: `internal/domain/repository/product.go`
- Modify: `internal/repository/product.go`

- [ ] **Step 1: Modify interface in domain repository**
  Update `Create` method signature in `internal/domain/repository/product.go` to add `createdBy uint64` parameter.
  ```go
  type ProductRepository interface{
  	Create(sku string, name string, description string, qty int, price float64, category []uint64, tag []uint64, createdBy uint64) (*entity.Product, error)
  	GetAll() ([]*entity.Product, error)
  	GetByID(id int) (*entity.Product, error)
  	Update(id int, sku string, name string, description string, qty int, price float64, category []uint64, tag []uint64) (*entity.Product, error)
  	Delete(id int) error
  }
  ```

- [ ] **Step 2: Modify Postgres implementation of repository**
  Update `Create` method in `internal/repository/product.go` to match the signature and save the `CreatedBy` field.
  ```go
  // Create
  func (r *ProductPostgresRepository) Create(sku string, name string, description string, qty int, price float64, category []uint64, tag []uint64, createdBy uint64) (*entity.Product, error) {
  	newItem := entity.Product{
  		SKU:         sku,
  		Name:        name,
  		Description: description,
  		Stock:       qty,
  		Price:       price,
  		Category:    category,
  		Tag:         tag,
  		CreatedBy:   createdBy,
  	}
  	err := r.db.Create(&newItem).Error
  	return &newItem, err
  }
  ```

- [ ] **Step 3: Verify Compilation**
  Run: `go build ./...`
  Expected: Success (except handlers calling `Create` which still need matching changes).

- [ ] **Step 4: Commit changes**
  Run: `git commit -am "repo: update Create signature to support CreatedBy"`

---

### Task 2: Fix RequirePermission Middleware

**Files:**
- Modify: `internal/delivery/http/middleware/auth.go`

- [ ] **Step 1: Update RequirePermission**
  Modify `RequirePermission` to check the `role` context key.
  ```go
  // Check Permission
  func RequirePermission(permissionRequired string) gin.HandlerFunc {
  	return func(c *gin.Context) {
  		roleVal, exists := c.Get("role")
  		if !exists {
  			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Role not found in context"})
  			c.Abort()
  			return
  		}
  		userRole, ok := roleVal.(string)
  		if !ok {
  			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Invalid role format in context"})
  			c.Abort()
  			return
  		}

  		// Check permission
  		roles := strings.Split(permissionRequired, ", ")
  		for _, role := range roles {
  			if strings.EqualFold(strings.TrimSpace(role), userRole) {
  				c.Next()
  				return
  			}
  		}

  		// If none of the allowed roles matched
  		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden", "message": "Insufficient permissions. Requires one of: " + permissionRequired})
  		c.Abort()
  	}
  }
  ```

- [ ] **Step 2: Verify Compilation**
  Run: `go build ./...`
  Expected: Success.

- [ ] **Step 3: Commit changes**
  Run: `git commit -am "middleware: fix RequirePermission to evaluate context role string"`

---

### Task 3: Enforce Ownership Check in Product Handlers

**Files:**
- Modify: `internal/delivery/http/product_handler.go`

- [ ] **Step 1: Pass Creator ID in CreateProduct**
  Extract `user_id` from Gin context, parse to `uint64`, and pass to `repo.Create`.
  ```go
  func (h *ProductHandler) CreateProduct(c *gin.Context) {
  	var payload entity.Product

  	if err := c.ShouldBindJSON(&payload); err != nil {
  		c.JSON(http.StatusBadRequest, gin.H{
  			"status":  "Bad Request",
  			"message": err.Error(),
  		})
  		return
  	}

  	// Get user_id from context
  	userIDStr, exists := c.Get("user_id")
  	if !exists {
  		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "User ID not found in context"})
  		return
  	}
  	userIDStrVal, ok := userIDStr.(string)
  	if !ok {
  		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Invalid User ID format"})
  		return
  	}
  	userID, err := strconv.ParseUint(userIDStrVal, 10, 64)
  	if err != nil {
  		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Invalid User ID format"})
  		return
  	}

  	newProduct, err := h.repo.Create(payload.SKU, payload.Name, payload.Description, payload.Stock, payload.Price, payload.Category, payload.Tag, userID)
  	if err != nil {
  		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": err.Error()})
  		return
  	}

  	c.JSON(http.StatusCreated, gin.H{
  		"status":   "Created",
  		"message":  "Product created successfully",
  		"products": newProduct,
  	})
  }
  ```

- [ ] **Step 2: Add Ownership Verification in UpdateProduct**
  Retrieve product first and check if role is `"seller"` that they match `CreatedBy`.
  ```go
  func (h *ProductHandler) UpdateProduct(c *gin.Context) {
  	idStr := c.Param("id")
  	id, err := strconv.Atoi(idStr)
  	if err != nil {
  		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid ID format"})
  		return
  	}
  	var payload entity.Product
  	if err := c.ShouldBindJSON(&payload); err != nil {
  		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": err.Error()})
  		return
  	}

  	// Get role from context
  	roleVal, exists := c.Get("role")
  	if !exists {
  		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Role not found in context"})
  		return
  	}
  	role, ok := roleVal.(string)
  	if !ok {
  		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Invalid role format"})
  		return
  	}

  	// Get product to check ownership
  	product, err := h.repo.GetByID(id)
  	if err != nil {
  		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": err.Error()})
  		return
  	}

  	if role == "seller" {
  		userIDStr, exists := c.Get("user_id")
  		if !exists {
  			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "User ID not found in context"})
  			return
  		}
  		userIDStrVal, ok := userIDStr.(string)
  		if !ok {
  			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Invalid User ID format"})
  			return
  		}
  		userID, err := strconv.ParseUint(userIDStrVal, 10, 64)
  		if err != nil {
  			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Invalid User ID format"})
  			return
  		}

  		if product.CreatedBy != userID {
  			c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden", "message": "You do not own this product"})
  			return
  		}
  	}

  	updated, err := h.repo.Update(id, payload.SKU, payload.Name, payload.Description, payload.Stock, payload.Price, payload.Category, payload.Tag)
  	if err != nil {
  		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": err.Error()})
  		return
  	}
  	c.JSON(http.StatusOK, gin.H{"status": "Ok", "message": "Product updated successfully", "products": updated})
  }
  ```

- [ ] **Step 3: Add Ownership Verification in DeleteProduct**
  Retrieve product first and check if role is `"seller"` that they match `CreatedBy`.
  ```go
  func (h *ProductHandler) DeleteProduct(c *gin.Context) {
  	idStr := c.Param("id")
  	id, err := strconv.Atoi(idStr)
  	if err != nil {
  		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid ID format"})
  		return
  	}

  	// Get role from context
  	roleVal, exists := c.Get("role")
  	if !exists {
  		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Role not found in context"})
  		return
  	}
  	role, ok := roleVal.(string)
  	if !ok {
  		c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Invalid role format"})
  		return
  	}

  	// Get product to check ownership
  	product, err := h.repo.GetByID(id)
  	if err != nil {
  		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": err.Error()})
  		return
  	}

  	if role == "seller" {
  		userIDStr, exists := c.Get("user_id")
  		if !exists {
  			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "User ID not found in context"})
  			return
  		}
  		userIDStrVal, ok := userIDStr.(string)
  		if !ok {
  			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Invalid User ID format"})
  			return
  		}
  		userID, err := strconv.ParseUint(userIDStrVal, 10, 64)
  		if err != nil {
  			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Invalid User ID format"})
  			return
  		}

  		if product.CreatedBy != userID {
  			c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden", "message": "You do not own this product"})
  			return
  		}
  	}

  	if err := h.repo.Delete(id); err != nil {
  		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": err.Error()})
  		return
  	}
  	c.JSON(http.StatusOK, gin.H{"status": "Ok", "message": "Product deleted successfully"})
  }
  ```

- [ ] **Step 4: Verify Compilation**
  Run: `go build ./...`
  Expected: Success.

- [ ] **Step 5: Commit changes**
  Run: `git commit -am "handlers: enforce product ownership and pass creator ID"`

---

## Verification Plan

### Automated Verification
Run `go build ./...` and `go test ./...` to ensure everything builds successfully and there are no lint issues.

### Manual Verification
1. Log in as a User with the role `seller` and create a product.
2. Verify that the product is saved with the correct `CreatedBy` matching the seller's `user_id`.
3. Try to update/delete that product using the same seller's token. (Should be allowed)
4. Try to update/delete that product using a different seller's token. (Should return `403 Forbidden`)
5. Try to update/delete that product using an `admin` token. (Should be allowed)
