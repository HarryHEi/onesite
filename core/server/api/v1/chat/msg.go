package chat

const (
	SystemMsgCode = iota
	UserMsgCode
)

// 表示一个消息
type WsMessage struct {
	Code int    `json:"code"`
	Src  string `json:"src"`
	Data string `json:"data"`
}
