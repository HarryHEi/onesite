package chat

import (
	"go.uber.org/zap"
	"gopkg.in/olahol/melody.v1"

	"onesite/common/log"
)

func Message(m *melody.Melody, session *melody.Session, bytes []byte) {
	err := m.BroadcastOthers(bytes, session)
	if err != nil {
		log.Error("BroadcastOthers failed", zap.Error(err))
		return
	}
}
