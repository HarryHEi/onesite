package redis

import (
	"context"
	"errors"
	"onesite/core/config"

	"github.com/go-redis/redis/v8"
)

var (
	_redis *redis.Client
	Cfg    = defaultConfig()
)

func InitRedis(options ...Option) error {
	for _, option := range options {
		option(Cfg)
	}
	_redis = redis.NewClient(&redis.Options{
		Addr:     Cfg.Addr,
		Password: Cfg.Password,
		DB:       Cfg.DB,
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

func defaultConfig() *config.RedisConfig {
	return &config.RedisConfig{
		Addr:     "172.172.177.191:6379",
		Password: "",
		DB:       0,
	}
}

type Option func(*config.RedisConfig)
