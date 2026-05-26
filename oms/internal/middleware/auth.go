package middleware

import (
	"net/http"
	"strings"

	"oms/internal/service"
	"oms/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "未提供认证令牌")
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "令牌格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]
		userID, username, err := authSvc.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "令牌无效或已过期")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", userID)
		c.Set("username", username)

		c.Next()
	}
}
