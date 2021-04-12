package dao

import (
	"errors"
	mongo2 "onesite/common/mongo"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	"onesite/common/orm"
	redis2 "onesite/common/redis"
	"onesite/core/model"
)

var (
	_dao *Dao
)

type Dao struct {
	Db    *gorm.DB
	Redis *redis.Client
	Mongo *mongo.Client
}

func InitDao() (err error) {
	if _dao != nil {
		return nil
	}

	err = orm.InitOrm()
	if err != nil {
		return err
	}

	db, err := orm.GetDb()
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	err = redis2.InitRedis()
	if err != nil {
		return err
	}

	rd, err := redis2.GetRedis()
	if err != nil {
		return err
	}

	err = mongo2.InitMongo()
	if err != nil {
		return err
	}

	mg, err := mongo2.GetMongo()
	if err != nil {
		return err
	}

	_dao = &Dao{
		db,
		rd,
		mg,
	}
	return nil
}

func GetDao() (*Dao, error) {
	if _dao == nil {
		return nil, errors.New("call InitDao before GetDao")
	}
	return _dao, nil
}
