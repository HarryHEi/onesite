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

func MakeUserMessage(name, msg string) []byte {
	msgObj := model.Message{
		Code: model.UserMsgCode,
		Src:  name,
		Data: msg,
	}
	msgStr, _ := json.Marshal(msgObj)
	return msgStr
}

func SaveUserMessage(name, msg string) error {
	return dao.SaveMessage(&model.Message{
		Code: model.UserMsgCode,
		Src:  name,
		Data: msg,
	})
}
