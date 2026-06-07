package auth

import "github.com/gin-gonic/gin"

// Auth returns a Gin middleware that performs unified API access permission checks.
// Currently a placeholder — implement the actual permission logic before production use.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO：统一校验接口访问权限
		c.Next()
	}
}
