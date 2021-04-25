package db

import (
	"context"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"onesite/core/config"
)

type MongoConfig struct {
	Uri string `toml:"uri"`
	Db  string `toml:"db"`
}

type MongoCli struct {
	Db *mongo.Database
}

func NewMongo() (*MongoCli, error) {
	var cfg MongoConfig
	_, err := toml.DecodeFile(config.GetCfgPath("mongo.toml"), &cfg)
	if err != nil {
		return nil, err
	}

	mongoUri, err := url.QueryUnescape(os.Getenv("MONGO_URI"))
	if err == nil && mongoUri != "" {
		cfg.Uri = mongoUri
	}
	mongoDb, err := url.QueryUnescape(os.Getenv("MONGO_DB"))
	if err == nil && mongoDb != "" {
		cfg.Db = mongoDb
	}

	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(cfg.Uri),
	)
	if err != nil {
		return nil, err
	}
	db := client.Database(cfg.Db)
	return &MongoCli{
		db,
	}, nil
}
