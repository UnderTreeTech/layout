package auth

import "github.com/gin-gonic/gin"

// Token returns a Gin middleware that performs unified token validity checks.
// Currently a placeholder — implement the actual token verification logic before production use.
func Token() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO：统一校验token有效性
		c.Next()
	}
}
