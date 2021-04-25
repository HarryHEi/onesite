package db

import (
	"errors"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"onesite/core/config"
)

type OrmConfig struct {
	DriverName  string `toml:"driver_name"`
	Dsn         string `toml:"dsn"`
	MaxOpenConn int    `toml:"max_open_conn"`
	MaxIdleConn int    `toml:"max_idle_conn"`
}

type OrmCli struct {
	Db *gorm.DB
}

func NewOrm() (*OrmCli, error) {
	var cfg OrmConfig
	_, err := toml.DecodeFile(config.GetCfgPath("db.toml"), &cfg)
	if err != nil {
		return nil, err
	}

	dbDsn, err := url.QueryUnescape(os.Getenv("DB_DSN"))
	if err == nil && dbDsn != "" {
		cfg.Dsn = dbDsn
	}

	var db *gorm.DB
	switch cfg.DriverName {
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{})
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(cfg.Dsn), &gorm.Config{})
	default:
		return nil, errors.New("unknown db driver: " + cfg.DriverName)
	}
	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDb.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDb.SetMaxIdleConns(cfg.MaxIdleConn)

	err = sqlDb.Ping()
	if err != nil {
		return nil, err
	}

	return &OrmCli{
		db,
	}, nil
}
