package middleware

import (
	"feishu/dep"

	"github.com/gin-gonic/gin"
)

const KeyDep = "key:dep"

func WithDependency(s *dep.HttpDependency) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(KeyDep, s)
		c.Next()
	}
}

func Dependency(c *gin.Context) *dep.HttpDependency {
	return c.MustGet(KeyDep).(*dep.HttpDependency)
}
