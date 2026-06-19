package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"study-golang-backend/pkg/auth"

	"github.com/gin-gonic/gin"
)

// Validate Bearer token
func AuthMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Check header has Authorization?
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Authorization header is required"})
			c.Abort()
			return
		}

		//Check Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Authorization header is required"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := auth.ValidateToken(parts[1], secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Invalid token"})
			c.Abort()
			return
		}

		// Injection clams into context
		c.Set("username", claims.Username)
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// Check Permission: Check permissionRequired this user has this role
func RequirePermission(permissionRequired string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Role from context
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "User role not found in context"})
			c.Abort()
			return
		}
		roleStr, ok := roleVal.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Invalid User role format"})
			c.Abort()
			return
		}

		// 1. Check role in permissionRequired
		allowed := false
		roles := strings.Split(permissionRequired, ", ")
		for _, role := range roles {
			if strings.EqualFold(roleStr, role) {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden", "message": "Insufficient permissions. Requires one of: " + permissionRequired})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Check Ownership: Check this permission need owner by getOwnerByID
func RequireOwnership(roleToCheckOwnership string, getOwnerByID func(id uint64) (string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Role from context
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "User role not found in context"})
			c.Abort()
			return
		}
		roleStr, ok := roleVal.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Invalid User role format"})
			c.Abort()
			return
		}

		// Check this permission need owner?
		listRolesOwnership := strings.Split(roleToCheckOwnership, ", ")
		for _, role := range listRolesOwnership {
			if strings.EqualFold(roleStr, strings.TrimSpace(role)) {
				// get id param http.../:id
				idStr := c.Param("id")
				id, err := strconv.Atoi(idStr)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid ID format"})
					c.Abort()
					return
				}
				// get id from query
				ownerID, err := getOwnerByID(uint64(id))
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "Failed to get owner"})
					c.Abort()
					return
				}
				// get user_id from context
				userIDVal, exists := c.Get("user_id")
				if !exists {
					c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "User ID not found in context"})
					c.Abort()
					return
				}
				userIDStr, ok := userIDVal.(string)
				if !ok {
					c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Invalid User ID format"})
					c.Abort()
					return
				}
				// check if user is owner
				if userIDStr != ownerID {
					c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden", "message": "Insufficient permissions. Requires ownership"})
					c.Abort()
					return
				}
			}
		}

		// not has in role that requier ownership
		c.Next()
	}
}
