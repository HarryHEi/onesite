package orm

import (
	"errors"
	"onesite/core/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	_db *gorm.DB
	Cfg = defaultConfig()
)

func InitOrm(options ...Option) (err error) {
	for _, option := range options {
		option(Cfg)
	}
	switch Cfg.DriverName {
	case "mysql":
		_db, err = gorm.Open(mysql.Open(Cfg.Dsn), &gorm.Config{})
	case "sqlite3":
		_db, err = gorm.Open(sqlite.Open(Cfg.Dsn), &gorm.Config{})
	default:
		return errors.New("unknown db driver: " + Cfg.DriverName)
	}
	if err != nil {
		return err
	}

	sqlDb, err := _db.DB()
	if err != nil {
		return err
	}

	sqlDb.SetMaxOpenConns(Cfg.MaxOpenConn)
	sqlDb.SetMaxIdleConns(Cfg.MaxIdleConn)
	return sqlDb.Ping()
}

func GetDb() (*gorm.DB, error) {
	if _db == nil {
		return nil, errors.New("call InitDB before GetDb")
	}

	return _db, nil
}

func defaultConfig() *config.DbConfig {
	return &config.DbConfig{
		DriverName:  "sqlite3",
		Dsn:         "sqlite3.db",
		MaxOpenConn: 10,
		MaxIdleConn: 5,
	}
}

type Option func(*config.DbConfig)

func DriverName(driverName string) Option {
	return func(config *config.DbConfig) {
		config.DriverName = driverName
	}
}

func Dsn(dsn string) Option {
	return func(config *config.DbConfig) {
		config.Dsn = dsn
	}
}

func MaxOpenConn(maxOpenConn int) Option {
	return func(config *config.DbConfig) {
		config.MaxOpenConn = maxOpenConn
	}
}

func MaxIdleConn(maxIdleConn int) Option {
	return func(config *config.DbConfig) {
		config.MaxIdleConn = maxIdleConn
	}
}
