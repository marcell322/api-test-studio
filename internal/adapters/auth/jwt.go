package auth

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint, secret string, expireHours int) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(userID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tkn.SignedString([]byte(secret))
}

// Middleware validates Bearer token and sets userID in context ("userID")
func Middleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "missing auth"}); return }
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid auth header"}); return
		}
		tknStr := parts[1]
		var claims jwt.RegisteredClaims
		_, err := jwt.ParseWithClaims(tknStr, &claims, func(token *jwt.Token) (interface{}, error) { return []byte(secret), nil })
		if err != nil { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid token"}); return }
		id, _ := strconv.Atoi(claims.Subject)
		c.Set("userID", uint(id))
		c.Next()
	}
}
