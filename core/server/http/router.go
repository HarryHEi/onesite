package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/olahol/melody.v1"

	"onesite/common/log"
	"onesite/core/middleware"
	"onesite/core/server/api/v1/auth"
	"onesite/core/server/api/v1/chat"
)

// middleware
func initMiddleware(s *Service) {
	err := middleware.InitAuthMiddleware()
	if err != nil {
		panic(fmt.Sprintf("InitAuthMiddleware failed. %v", err))
	}
	s.S.Use(middleware.Logger(), gin.Recovery())
}

func initApiV1(s *Service) {
	v1Router := s.S.Group("/api/v1")

	authMiddleware := middleware.GetAuthMiddleware()
	authRouter := v1Router.Group("/auth")
	{
		authRouter.POST("/login", authMiddleware.LoginHandler)
		authRouter.GET("/refresh", authMiddleware.RefreshHandler)

		authRouter.GET("/user/info", authMiddleware.MiddlewareFunc(), auth.UserInfo())
	}
}

func initWsV1(s *Service) {
	wsRouter := s.S.Group("/ws")

	// 通过query params认证
	authMiddleware := middleware.GetAuthMiddleware()
	wsRouter.GET(
		"/v1/chat",
		authMiddleware.MiddlewareFunc(),
		middleware.ParseUserMiddleware(),
		func(c *gin.Context) {
			user := middleware.ParseUser(c)
			err := s.M.HandleRequestWithKeys(
				c.Writer,
				c.Request,
				map[string]interface{}{
					"user": user,
				},
			)
			if err != nil {
				log.Error("melody HandleRequest failed", zap.Error(err))
			}
		})

	// 连接建立
	s.M.HandleConnect(func(session *melody.Session) {
		chat.Login(session)
	})

	// 消息
	s.M.HandleMessage(func(session *melody.Session, bytes []byte) {
		chat.Message(s.M, session, bytes)
	})

	// 连接断开
	s.M.HandleClose(func(session *melody.Session, i int, s string) error {
		chat.Logout(session)
		return nil
	})
}

func initBasicRouter(s *Service) {
	s.S.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}

func InitRouter(s *Service) {
	initMiddleware(s)
	initBasicRouter(s)
	initApiV1(s)
	initWsV1(s)
}
