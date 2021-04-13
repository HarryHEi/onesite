package chat

import (
	"encoding/json"

	"onesite/core/dao"
	"onesite/core/model"
)

func MakeSystemMessage(msg string) []byte {
	msgObj := model.Message{
		Code: model.SystemMsgCode,
		Src:  "",
		Data: msg,
	}
	msgStr, _ := json.Marshal(msgObj)
	return msgStr
}

func MakeUserMessage(username, name, msg string) []byte {
	msgObj := model.Message{
		Code: model.UserMsgCode,
		Src:  name + "(" + username + ")",
		Data: msg,
	}
	msgStr, _ := json.Marshal(msgObj)
	return msgStr
}

func SaveUserMessage(username, name, msg string) error {
	return dao.SaveMessage(&model.Message{
		Code: model.UserMsgCode,
		Src:  name + "(" + username + ")",
		Data: msg,
	})
}
