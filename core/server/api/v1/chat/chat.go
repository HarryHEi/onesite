package chat

import (
	"fmt"

	"go.uber.org/zap"
	"gopkg.in/olahol/melody.v1"

	"onesite/common/log"
	"onesite/core/middleware"
)

// 登录
func Login(session *melody.Session) {
	userInstance, ok := middleware.ParseWsUser(session)
	if !ok {
		log.Error("unauthorized")
		_ = session.Close()
		return
	}

	fmt.Println(userInstance.Username, userInstance.Name, "Login")
}

// 消息
func Message(m *melody.Melody, session *melody.Session, bytes []byte) {
	userInstance, ok := middleware.ParseWsUser(session)
	if !ok {
		log.Error("unauthorized")
		_ = session.Close()
		return
	}

	fmt.Println(userInstance.Username, userInstance.Name, string(bytes))

	err := m.BroadcastOthers(bytes, session)
	if err != nil {
		log.Error("BroadcastOthers failed", zap.Error(err))
		_ = session.Close()
		return
	}
}

// 登出
func Logout(session *melody.Session) {
	userInstance, ok := middleware.ParseWsUser(session)
	if !ok {
		log.Error("unauthorized")
		_ = session.Close()
		return
	}

	fmt.Println(userInstance.Username, userInstance.Name, "Logout")
}
