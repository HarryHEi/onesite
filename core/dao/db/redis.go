package db

import (
	"context"
	"errors"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis/v8"

	"onesite/core/config"
)

type RedisConfig struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

type RedisCli struct {
	Db *redis.Client
}

func NewRedis() (*RedisCli, error) {
	var cfg RedisConfig
	_, err := toml.DecodeFile(config.GetCfgPath("redis.toml"), &cfg)
	if err != nil {
		return nil, err
	}

	redisAddr, err := url.QueryUnescape(os.Getenv("REDIS_ADDR"))
	if err == nil && redisAddr != "" {
		cfg.Addr = redisAddr
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	v := client.Ping(context.Background())
	if v.Val() != "PONG" {
		return nil, errors.New("ping redis failed")
	}
	return &RedisCli{
		client,
	}, nil
}
