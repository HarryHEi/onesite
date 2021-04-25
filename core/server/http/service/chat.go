package service

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/olahol/melody.v1"

	"onesite/core/log"
	"onesite/core/middleware"
	"onesite/core/model"
	"onesite/core/server/http/rest"
)

func (s *Service) MakeSystemMessage(msg string) []byte {
	msgObj := model.Message{
		Code: model.SystemMsgCode,
		Src:  "",
		Data: msg,
	}
	msgStr, _ := json.Marshal(msgObj)
	return msgStr
}

func (s *Service) MakeUserMessage(name, msg string) []byte {
	msgObj := model.Message{
		Code: model.UserMsgCode,
		Src:  name,
		Data: msg,
	}
	msgStr, _ := json.Marshal(msgObj)
	return msgStr
}

func (s *Service) SaveUserMessage(name, msg string) error {
	return s.Dao.SaveMessage(&model.Message{
		Code: model.UserMsgCode,
		Src:  name,
		Data: msg,
	})
}

// Login 登录
func (s *Service) Login(m *melody.Melody, session *melody.Session) {
	userInstance, ok := middleware.ParseWsUser(session)
	if !ok {
		log.Error("unauthorized")
		_ = session.Close()
		return
	}

	_ = m.Broadcast(s.MakeSystemMessage(userInstance.Name + " Login."))
}

// Message 消息
func (s *Service) Message(m *melody.Melody, session *melody.Session, bytes []byte) {
	userInstance, ok := middleware.ParseWsUser(session)
	if !ok {
		log.Error("unauthorized")
		_ = session.Close()
		return
	}

	go func() {
		err := s.SaveUserMessage(
			userInstance.Name,
			string(bytes),
		)
		if err != nil {
			log.Error("SaveUserMessage Failed", zap.Error(err))
		}
	}()
	err := m.BroadcastOthers(
		s.MakeUserMessage(
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
func (s *Service) Logout(m *melody.Melody, session *melody.Session) {
	userInstance, ok := middleware.ParseWsUser(session)
	if !ok {
		log.Error("unauthorized")
		_ = session.Close()
		return
	}

	err := m.Broadcast(s.MakeSystemMessage(userInstance.Name + " Logout."))
	if err != nil {
		log.Error("Broadcast failed", zap.Error(err))
		_ = session.Close()
		return
	}
}

// MessageHistory 历史消息
func (s *Service) MessageHistory() func(c *gin.Context) {
	return func(c *gin.Context) {
		messages, err := s.Dao.QueryMessage(20)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		rest.Success(c, messages)
	}
}
