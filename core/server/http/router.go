package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/olahol/melody.v1"

	"onesite/common/log"
	"onesite/core/middleware"
	"onesite/core/server/api/v1/admin"
	"onesite/core/server/api/v1/chat"
	"onesite/core/server/api/v1/user"
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

	// 认证
	authRouter := v1Router.Group("/auth")
	{
		authRouter.POST("/login", authMiddleware.LoginHandler)
		authRouter.GET("/refresh", authMiddleware.RefreshHandler)
	}

	// 用户信息
	userRouter := v1Router.Group("/user")
	userRouter.Use(
		authMiddleware.MiddlewareFunc(),
		middleware.ParseUserMiddleware(),
	)
	{
		userRouter.GET("/info", user.Info())
	}

	// 管理员
	adminRouter := v1Router.Group("/admin")
	adminRouter.Use(
		authMiddleware.MiddlewareFunc(),
		middleware.ParseUserMiddleware(),
		middleware.AdminPermissionMiddleware(),
	)
	{
		adminRouter.GET("/users", admin.ListUsers())
		adminRouter.POST("/user", admin.CreateUser())
		adminRouter.DELETE("/user/:pk", admin.DeleteUser())
		adminRouter.PATCH("/user/:pk", admin.PatchUpdateUser())
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
			u := middleware.ParseUser(c)
			err := s.M.HandleRequestWithKeys(
				c.Writer,
				c.Request,
				map[string]interface{}{
					"user": u,
				},
			)
			if err != nil {
				log.Error("melody HandleRequest failed", zap.Error(err))
			}
		})

	// 连接建立
	s.M.HandleConnect(func(session *melody.Session) {
		chat.Login(s.M, session)
	})

	// 消息
	s.M.HandleMessage(func(session *melody.Session, bytes []byte) {
		chat.Message(s.M, session, bytes)
	})

	// 连接断开
	s.M.HandleClose(func(session *melody.Session, _ int, _ string) error {
		chat.Logout(s.M, session)
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
