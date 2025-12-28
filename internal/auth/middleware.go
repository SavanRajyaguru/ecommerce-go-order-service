package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/savanrajyaguru/ecommerce-go-order-service/pkg/logger"
	"github.com/savanrajyaguru/ecommerce-go-order-service/pkg/utils"
	"go.uber.org/zap"
)

const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.JSONError(c, 401, "Authorization header required", "missing_token")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.JSONError(c, 401, "Invalid authorization header format", "invalid_token")
			c.Abort()
			return
		}

		claims, err := ValidateToken(parts[1])
		if err != nil {
			logger.Log.Warn("Invalid token", zap.Error(err))
			utils.JSONError(c, 401, "Invalid or expired token", "unauthorized")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.JSONError(c, 403, "Role not found in context", "forbidden")
			c.Abort()
			return
		}

		if role.(string) != requiredRole && role.(string) != RoleAdmin {
			// Admins can access everything? Or strict?
			// Prompt: "ADMIN -> view all orders". "USER -> create/view own".
			utils.JSONError(c, 403, "Insufficient permissions", "forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}
