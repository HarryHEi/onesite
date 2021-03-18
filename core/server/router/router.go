package router

import (
	"github.com/gin-gonic/gin"
	"onesite/core/middleware"
	"onesite/core/server/api/v1/auth"
)

// auth
func initAuthRouter(_s *gin.RouterGroup) {
	authRouter := _s.Group("/auth")
	{
		authRouter.POST("/login", auth.Login())
	}
}

// middleware
func initMiddleware(_s *gin.Engine) {
	_s.Use(middleware.Logger(), gin.Recovery())
}

func initApiV1(_s *gin.Engine) {
	v1Router := _s.Group("/api/v1")
	initAuthRouter(v1Router)
}

func initBasicRouter(_s *gin.Engine) {
	_s.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}

func InitRouter(_s *gin.Engine) {
	initMiddleware(_s)
	initBasicRouter(_s)
	initApiV1(_s)
}
