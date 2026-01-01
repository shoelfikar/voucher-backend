package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shoelfikar/voucher-management-system/internal/delivery/http/response"
	"github.com/shoelfikar/voucher-management-system/pkg/jwt"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(jwtService jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("Missing authorization header"))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// No "Bearer " prefix found
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("Invalid authorization header format"))
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("Invalid or expired token"))
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}
