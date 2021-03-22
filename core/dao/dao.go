package dao

import (
	"errors"

	redis2 "onesite/common/redis"
	"onesite/core/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"onesite/common/orm"
	"onesite/core/config"
)

var (
	_dao *Dao
)

type Dao struct {
	Db    *gorm.DB
	Redis *redis.Client
}

func InitDao() (err error) {
	if _dao != nil {
		return nil
	}

	err = orm.InitOrm(
		orm.DriverName(config.CoreCfg.Db.DriverName),
		orm.Dsn(config.CoreCfg.Db.Dsn),
		orm.MaxOpenConn(config.CoreCfg.Db.MaxOpenConn),
		orm.MaxIdleConn(config.CoreCfg.Db.MaxIdleConn),
	)
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

	_dao = &Dao{
		db,
		rd,
	}
	return nil
}

func GetDao() (*Dao, error) {
	if _dao == nil {
		return nil, errors.New("call InitDao before GetDao")
	}
	return _dao, nil
}
