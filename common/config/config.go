package config

import (
	"net/url"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	CoreCfg = defaultConfig()
)

func Init(configFile string) error {
	if configFile == "" {
		return nil
	}

	_, err := toml.DecodeFile(configFile, &CoreCfg)
	if err != nil {
		return err
	}

	updateFromEnv()

	return err
}

// 从环境变量获取配置
func updateFromEnv() {
	// Db -> Dsn
	dbDsn, err := url.QueryUnescape(os.Getenv("DB_DSN"))
	if err == nil && dbDsn != "" {
		CoreCfg.Db.Dsn = dbDsn
	}

	// Redis -> Addr
	redisAddr, err := url.QueryUnescape(os.Getenv("REDIS_ADDR"))
	if err == nil && redisAddr != "" {
		CoreCfg.Redis.Addr = redisAddr
	}

	// Mongo -> Uri/Db
	mongoUri, err := url.QueryUnescape(os.Getenv("MONGO_URI"))
	if err == nil && redisAddr != "" {
		CoreCfg.Mongo.Uri = mongoUri
	}
	mongoDb, err := url.QueryUnescape(os.Getenv("MONGO_DB"))
	if err == nil && mongoDb != "" {
		CoreCfg.Mongo.Db = mongoDb
	}
}

func defaultConfig() *CoreConfig {
	return &CoreConfig{
		SecretKey: "A long string.",
		Server: ServerConfig{
			"0.0.0.0",
			8000,
		},
		Db: DbConfig{
			"mysql",
			"herui:Admin@123@tcp(172.172.177.191:3306)/onesite_dev?charset=utf8mb4&parseTime=true",
			10,
			5,
		},
		Redis: RedisConfig{
			"172.172.177.191:6379",
			"",
			0,
		},
		Auth: AuthConfig{
			duration{
				time.Hour,
			},
		},
		Mongo: MongoConfig{
			"mongodb://172.172.177.191:27017/",
			"onesite",
		},
		Weed: WeedConfig{
			"http",
			"172.172.177.191",
			9333,
			"172.172.177.191",
			8080,
		},
	}
}

type CoreConfig struct {
	SecretKey string       `toml:"secret_key"`
	Server    ServerConfig `toml:"server"`
	Db        DbConfig     `toml:"db"`
	Redis     RedisConfig  `toml:"redis"`
	Auth      AuthConfig   `toml:"auth"`
	Mongo     MongoConfig  `toml:"mongo"`
	Weed      WeedConfig   `toml:"weed"`
}

type ServerConfig struct {
	Bind string `toml:"bind"`
	Port int    `toml:"port"`
}

type DbConfig struct {
	DriverName  string `toml:"driver_name"`
	Dsn         string `toml:"dsn"`
	MaxOpenConn int    `toml:"max_open_conn"`
	MaxIdleConn int    `toml:"max_idle_conn"`
}

type RedisConfig struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

type MongoConfig struct {
	Uri string `toml:"uri"`
	Db  string `toml:"db"`
}

type WeedConfig struct {
	Protocol   string `toml:"protocol"`
	VolumeHost string `toml:"volume_host"`
	VolumePort int    `toml:"volume_port"`
	FsHost     string `toml:"fs_host"`
	FsPort     int    `toml:"fs_port"`
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type AuthConfig struct {
	Timeout duration `toml:"timeout"`
}
