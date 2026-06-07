package http

import (
	"github.com/gin-gonic/gin"
)

// registerAPI registers all HTTP API routes onto the given gin.Engine.
// Routes are grouped under the /api prefix.
func registerAPI(engine *gin.Engine) {
	login := engine.Group("/api")
	{
		login.GET("/app/user", getUserInfo) // 查询用户信息
	}
}
