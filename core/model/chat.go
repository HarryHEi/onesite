package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	SystemMsgCode = iota
	UserMsgCode
)

const (
	MessageCollectionName = "chat.message"
)

type Message struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Code     int                `json:"code" bson:"code"`
	Src      string             `json:"src" bson:"src"`
	Data     string             `json:"data" bson:"data"`
	Creation time.Time          `json:"creation" bson:"creation,omitempty"`
}
