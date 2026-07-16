package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/marcell322/api-test-studio/internal/adapters/auth"
)

// AuthMiddleware validates JWT from Authorization header and stores userID in context.
// Expects: Authorization: Bearer <token>
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// read Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "missing authorization header",
			})
			return
		}

		// parse Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "invalid authorization header format",
			})
			return
		}

		tokenStr := parts[1]

		// validate token
		userID, err := auth.ValidateToken(tokenStr, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "invalid token",
			})
			return
		}

		// store userID in context
		c.Set("userID", userID)
		c.Next()
	}
}
