package chat

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/olahol/melody.v1"

	"onesite/common/log"
	"onesite/common/rest"
	"onesite/core/dao"
	"onesite/core/middleware"
)

// Login 登录
func Login(m *melody.Melody, session *melody.Session) {
	userInstance, ok := middleware.ParseWsUser(session)
	if !ok {
		log.Error("unauthorized")
		_ = session.Close()
		return
	}

	_ = m.Broadcast(MakeSystemMessage(userInstance.Name + " Login."))
}

// Message 消息
func Message(m *melody.Melody, session *melody.Session, bytes []byte) {
	userInstance, ok := middleware.ParseWsUser(session)
	if !ok {
		log.Error("unauthorized")
		_ = session.Close()
		return
	}

	go func() {
		err := SaveUserMessage(
			userInstance.Name,
			string(bytes),
		)
		if err != nil {
			log.Error("SaveUserMessage Failed", zap.Error(err))
		}
	}()
	err := m.BroadcastOthers(
		MakeUserMessage(
			userInstance.Name,
			string(bytes),
		),
		session,
	)
	if err != nil {
		log.Error("BroadcastOthers failed", zap.Error(err))
		_ = session.Close()
		return
	}
}

// Logout 登出
func Logout(m *melody.Melody, session *melody.Session) {
	userInstance, ok := middleware.ParseWsUser(session)
	if !ok {
		log.Error("unauthorized")
		_ = session.Close()
		return
	}

	err := m.Broadcast(MakeSystemMessage(userInstance.Name + " Logout."))
	if err != nil {
		log.Error("Broadcast failed", zap.Error(err))
		_ = session.Close()
		return
	}
}

// MessageHistory 历史消息
func MessageHistory() func(c *gin.Context) {
	return func(c *gin.Context) {
		messages, err := dao.QueryMessage(20)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		rest.Success(c, messages)
	}
}
