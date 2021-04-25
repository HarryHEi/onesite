package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/olahol/melody.v1"

	"onesite/core/log"
	"onesite/core/middleware"
)

func InitRouter(s *Server) error {
	s.S.Use(middleware.Logger(), gin.Recovery())

	s.S.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	limiter := middleware.NewRateLimiter(500)

	v1Router := s.S.Group("/api/v1")
	v1Router.Use(limiter.Middleware())

	authMiddleware, err := middleware.NewAuthMiddleware(s.Dao)
	if err != nil {
		return err
	}

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
		middleware.ParseUserMiddleware(s.Dao),
	)
	{
		userRouter.GET("/info", s.Svc.Info())
		userRouter.POST("/avatar", s.Svc.UploadAvatar())
		userRouter.POST("/password", s.Svc.UpdatePassword())
	}

	// 管理员
	adminRouter := v1Router.Group("/admin")
	adminRouter.Use(
		authMiddleware.MiddlewareFunc(),
		middleware.ParseUserMiddleware(s.Dao),
		middleware.AdminPermissionMiddleware(),
	)
	{
		adminRouter.GET("/users", s.Svc.ListUsers())
		adminRouter.POST("/user", s.Svc.CreateUser())
		adminRouter.DELETE("/user/:pk", s.Svc.DeleteUser())
		adminRouter.PATCH("/user/:pk", s.Svc.PatchUpdateUser())
	}

	// chat
	chatRouter := v1Router.Group("/chat")
	chatRouter.Use(
		authMiddleware.MiddlewareFunc(),
	)
	{
		chatRouter.GET("/history", s.Svc.MessageHistory())
	}

	// fs
	fsRouter := v1Router.Group("/fs")
	fsRouter.Use(
		authMiddleware.MiddlewareFunc(),
		middleware.ParseUserMiddleware(s.Dao),
		middleware.StrangerDeniedMiddleware(),
	)
	{
		fsRouter.GET("/list", s.Svc.ListFiles())
		fsRouter.POST("/upload", s.Svc.UploadFile())
		fsRouter.GET("/download/:pk", s.Svc.DownloadFile())
		fsRouter.DELETE("/delete/:pk", s.Svc.DeleteFile())
		fsRouter.POST("/export/:pk", s.Svc.SetExportFile())
	}

	// export
	exportRouter := v1Router.Group("/export")
	{
		exportRouter.GET("/fs/:pk", s.Svc.ExportFile())
		exportRouter.GET("/avatar/:pk", s.Svc.ExportAvatar())
	}

	wsRouter := s.S.Group("/ws")

	wsRouter.GET(
		"/v1/chat",
		authMiddleware.MiddlewareFunc(),
		middleware.ParseUserMiddleware(s.Dao),
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
		s.Svc.Login(s.M, session)
	})

	// 消息
	s.M.HandleMessage(func(session *melody.Session, bytes []byte) {
		// 消息长度限制
		if len(bytes) > 256 {
			bytes = bytes[:256]
		}

		s.Svc.Message(s.M, session, bytes)
	})

	// 连接断开
	s.M.HandleClose(func(session *melody.Session, _ int, _ string) error {
		s.Svc.Logout(s.M, session)
		return nil
	})

	return nil
}
