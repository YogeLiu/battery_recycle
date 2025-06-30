package v1

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService *services.AuthService
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware checks for valid JWT token
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, &models.Response{
				Code: models.CodeUnauthorized,
				Msg:  "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusOK, &models.Response{
				Code: models.CodeUnauthorized,
				Msg:  "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]
		user, err := m.authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusOK, &models.Response{
				Code: models.CodeUnauthorized,
				Msg:  "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Store user in context
		c.Set("user", user)
		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusOK, &models.Response{
				Code: models.CodeUnauthorized,
				Msg:  "User not authenticated",
			})
			c.Abort()
			return
		}

		userModel := user.(*models.User)
		if userModel.Role != role {
			c.JSON(http.StatusOK, &models.Response{
				Code: models.CodeForbidden,
				Msg:  "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
