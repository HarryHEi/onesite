package orm

import (
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"onesite/common/config"
)

var (
	_db *gorm.DB
)

func InitOrm(options ...Option) (err error) {
	for _, option := range options {
		option(&config.CoreCfg.Db)
	}
	switch config.CoreCfg.Db.DriverName {
	case "mysql":
		_db, err = gorm.Open(mysql.Open(config.CoreCfg.Db.Dsn), &gorm.Config{})
	case "sqlite3":
		_db, err = gorm.Open(sqlite.Open(config.CoreCfg.Db.Dsn), &gorm.Config{})
	default:
		return errors.New("unknown db driver: " + config.CoreCfg.Db.DriverName)
	}
	if err != nil {
		return err
	}

	sqlDb, err := _db.DB()
	if err != nil {
		return err
	}

	sqlDb.SetMaxOpenConns(config.CoreCfg.Db.MaxOpenConn)
	sqlDb.SetMaxIdleConns(config.CoreCfg.Db.MaxIdleConn)
	return sqlDb.Ping()
}

func GetDb() (*gorm.DB, error) {
	if _db == nil {
		return nil, errors.New("call InitDB before GetDb")
	}

	return _db, nil
}

type Option func(*config.DbConfig)

//func DriverName(driverName string) Option {
//	return func(config *config.DbConfig) {
//		config.DriverName = driverName
//	}
//}
//
//func Dsn(dsn string) Option {
//	return func(config *config.DbConfig) {
//		config.Dsn = dsn
//	}
//}
//
//func MaxOpenConn(maxOpenConn int) Option {
//	return func(config *config.DbConfig) {
//		config.MaxOpenConn = maxOpenConn
//	}
//}
//
//func MaxIdleConn(maxIdleConn int) Option {
//	return func(config *config.DbConfig) {
//		config.MaxIdleConn = maxIdleConn
//	}
//}
