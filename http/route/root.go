package route

import (
	"github.com/gin-gonic/gin"
)

func Root(opts ...func(engine *gin.Engine)) func(s *gin.Engine) {
	return func(s *gin.Engine) {
		for _, opt := range opts {
			opt(s)
		}

		// common middleware
		s.Use(
			func(c *gin.Context) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
				c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
				if c.Request.Method == "OPTIONS" {
					c.AbortWithStatus(204)
				}
			},
		)

	}
}
