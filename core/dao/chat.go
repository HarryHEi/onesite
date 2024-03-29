package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"onesite/core/model"
)

func (dao *Dao) SaveMessage(message *model.Message) error {
	collection := dao.Mongo.Db.Collection(model.MessageCollectionName)
	message.Creation = time.Now()
	_, err := collection.InsertOne(context.Background(), message)
	if err != nil {
		return err
	}
	return nil
}

func (dao *Dao) QueryMessage(limit int64) ([]model.Message, error) {
	if limit > 20 {
		limit = 20
	}

	collection := dao.Mongo.Db.Collection(model.MessageCollectionName)

	opts := options.FindOptions{
		Limit: &limit,
		Sort:  bson.D{{"creation", -1}},
	}
	cur, err := collection.Find(context.Background(), bson.D{}, &opts)
	if err != nil {
		return nil, err
	}
	var messages []model.Message
	err = cur.All(context.Background(), &messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
