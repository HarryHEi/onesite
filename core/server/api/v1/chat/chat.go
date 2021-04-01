package chat

import (
	"encoding/json"

	"go.uber.org/zap"
	"gopkg.in/olahol/melody.v1"

	"onesite/common/log"
	"onesite/core/middleware"
)

// 登录
func Login(m *melody.Melody, session *melody.Session) {
	userInstance, ok := middleware.ParseWsUser(session)
	if !ok {
		log.Error("unauthorized")
		_ = session.Close()
		return
	}

	greetingMsg := WsMessage{
		SystemMsgCode,
		"",
		userInstance.Name + " Login.",
	}
	greetingMsgStr, _ := json.Marshal(greetingMsg)

	_ = m.Broadcast(greetingMsgStr)
}

// 消息
func Message(m *melody.Melody, session *melody.Session, bytes []byte) {
	userInstance, ok := middleware.ParseWsUser(session)
	if !ok {
		log.Error("unauthorized")
		_ = session.Close()
		return
	}

	userMessage := WsMessage{
		UserMsgCode,
		userInstance.Name + "(" + userInstance.Username + ")",
		string(bytes),
	}
	userMessageStr, _ := json.Marshal(userMessage)

	err := m.BroadcastOthers(userMessageStr, session)
	if err != nil {
		log.Error("BroadcastOthers failed", zap.Error(err))
		_ = session.Close()
		return
	}
}

// 登出
func Logout(m *melody.Melody, session *melody.Session) {
	userInstance, ok := middleware.ParseWsUser(session)
	if !ok {
		log.Error("unauthorized")
		_ = session.Close()
		return
	}

	byeMsg := WsMessage{
		SystemMsgCode,
		"",
		userInstance.Name + " Logout.",
	}
	byeMsgStr, _ := json.Marshal(byeMsg)

	err := m.Broadcast(byeMsgStr)
	if err != nil {
		log.Error("Broadcast failed", zap.Error(err))
		_ = session.Close()
		return
	}
}
