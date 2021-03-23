package middleware

import (
	"gopkg.in/olahol/melody.v1"

	"onesite/core/model"
)

// 从会话keys解析用户实例
func ParseWsUser(session *melody.Session) (*model.User, bool) {
	user, exist := session.Get("user")
	if !exist {
		return nil, false
	}
	userInstance, ok := user.(*model.User)
	return userInstance, ok
}
