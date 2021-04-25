package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"onesite/core/dao/db"
	"onesite/core/dao/external"
	"onesite/core/dao/sync"
	"onesite/core/model"
)

type Dao struct {
	Orm    *db.OrmCli
	Redis  *db.RedisCli
	Mongo  *db.MongoCli
	RDLock *sync.RDLock
	Weed   *external.Weed
}

func NewDao() (*Dao, error) {
	orm, err := db.NewOrm()
	if err != nil {
		return nil, err
	}

	rd, err := db.NewRedis()
	if err != nil {
		return nil, err
	}

	mg, err := db.NewMongo()
	if err != nil {
		return nil, err
	}

	rdl := sync.NewRDLock(rd.Db)

	wd, err := external.NewWeed()
	if err != nil {
		return nil, err
	}

	return &Dao{
		orm,
		rd,
		mg,
		rdl,
		wd,
	}, nil
}

func (dao *Dao) Migrate() error {
	// DB
	err := dao.Orm.Db.AutoMigrate(&model.User{}, &model.File{})
	if err != nil {
		return err
	}

	// Mongo
	err = dao.Mongo.Db.CreateCollection(context.Background(), model.MessageCollectionName)
	if err != nil {
		return err
	}
	indexes := dao.Mongo.Db.Collection(model.MessageCollectionName).Indexes()
	_, err = indexes.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{
				{"creation", 1},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}
