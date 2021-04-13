package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"onesite/common/config"
)

var (
	_mongo *mongo.Database
)

func InitMongo(ops ...Option) (err error) {
	for _, option := range ops {
		option(&config.CoreCfg.Mongo)
	}
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(config.CoreCfg.Mongo.Uri),
	)
	if err != nil {
		return err
	}
	_mongo = client.Database(config.CoreCfg.Mongo.Db)
	return err
}

func GetMongo() (*mongo.Database, error) {
	if _mongo == nil {
		return nil, errors.New("call InitMongo before GetMongo")
	}
	return _mongo, nil
}

type Option func(config *config.MongoConfig)

//func Uri(uri string) Option {
//	return func(config *config.MongoConfig) {
//		config.Uri = uri
//	}
//}
//
//func Db(db string) Option {
//	return func(config *config.MongoConfig) {
//		config.Db = db
//	}
//}
