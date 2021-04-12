package redis

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"

	"onesite/common/config"
)

var (
	_redis *redis.Client
)

func InitRedis(options ...Option) error {
	for _, option := range options {
		option(&config.CoreCfg.Redis)
	}
	_redis = redis.NewClient(&redis.Options{
		Addr:     config.CoreCfg.Redis.Addr,
		Password: config.CoreCfg.Redis.Password,
		DB:       config.CoreCfg.Redis.DB,
	})
	v := _redis.Ping(context.Background())
	if v.Val() != "PONG" {
		return errors.New("ping redis failed")
	}
	return nil
}

func GetRedis() (*redis.Client, error) {
	if _redis == nil {
		return nil, errors.New("call InitRedis before GetRedis")
	}
	return _redis, nil
}

type Option func(*config.RedisConfig)

func Addr(addr string) Option {
	return func(config *config.RedisConfig) {
		config.Addr = addr
	}
}

func Password(password string) Option {
	return func(config *config.RedisConfig) {
		config.Password = password
	}
}

func Db(db int) Option {
	return func(config *config.RedisConfig) {
		config.DB = db
	}
}
