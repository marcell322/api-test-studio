package auth

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// GetSecretFromEnv reads the JWT secret from JWT_SECRET env var.
// Returns empty string if not set; callers should handle that.
func GetSecretFromEnv() string {
	return os.Getenv("JWT_SECRET")
}

// GenerateToken creates a signed JWT for the given user ID and expiry hours.
func GenerateToken(userID uint, secret string, expireHours int) (string, error) {
	if secret == "" {
		return "", fmt.Errorf("jwt secret is empty")
	}
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(userID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tkn.SignedString([]byte(secret))
}

// ValidateToken parses and validates a token string using the provided secret.
// Returns the user ID (from subject) if valid.
func ValidateToken(tokenStr, secret string) (uint, error) {
	if secret == "" {
		return 0, fmt.Errorf("jwt secret is empty")
	}
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		// ensure signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}
	if claims.Subject == "" {
		return 0, fmt.Errorf("token missing subject")
	}
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, fmt.Errorf("invalid subject in token: %w", err)
	}
	return uint(id), nil
}

// Middleware validates Bearer token and sets userID in context ("userID").
func Middleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "missing auth"})
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid auth header"})
			return
		}
		tknStr := parts[1]
		id, err := ValidateToken(tknStr, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid token"})
			return
		}
		c.Set("userID", id)
		c.Next()
	}
}
