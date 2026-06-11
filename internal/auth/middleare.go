package auth

import (
	"net/http"
	"strings"

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
		claims, err := ValidateToken(parts[1], secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized", "message": "Invalid token"})
			c.Abort()
			return
		}

		// Injection clams
		c.Set("username", claims.Username)
		c.Set("user_id", claims.UserID)
		c.Set("can_read", claims.CanRead)
		c.Set("can_write", claims.CanWrite)
		c.Set("can_update", claims.CanUpdate)
		c.Set("can_delete", claims.CanDelete)
		c.Next()
	}
}

// Check Permission
func RequirePermission(permissionKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read permission key
		val, exists := c.Get(permissionKey)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden", "message": "Required permission not found"})
			c.Abort()
			return
		}

		// Check permission
		allowed, ok := val.(bool)
		if !ok || !allowed {
			c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden", "message": "Insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}