package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth ...
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func Token(c *gin.Context) string {
	value := c.GetHeader("Authorization")
	if value != "" && strings.HasPrefix(value, "Bearer ") {
		return strings.TrimPrefix(value, "Bearer ")
	}
	return ""
}
