package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"onesite/core/model"
)

func SaveMessage(message *model.Message) error {
	dao, err := GetDao()
	if err != nil {
		return err
	}
	collection := dao.Mongo.Collection(model.MessageCollectionName)
	message.Creation = time.Now()
	_, err = collection.InsertOne(context.Background(), message)
	if err != nil {
		return err
	}
	return nil
}

func QueryMessage(limit int64) ([]model.Message, error) {
	dao, err := GetDao()
	if err != nil {
		return nil, err
	}
	if limit > 20 {
		limit = 20
	}

	collection := dao.Mongo.Collection(model.MessageCollectionName)

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
