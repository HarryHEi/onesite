package redis_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"onesite/common/redis"
)

func TestInitRedis(t *testing.T) {
	_, err := redis.GetRedis()
	require.NotNil(t, err)
	err = redis.InitRedis(
		redis.Addr("172.172.177.191:6379"),
		redis.Password(""),
		redis.Db(0),
	)
	require.Nil(t, err)
	cli, _ := redis.GetRedis()
	require.NotNil(t, cli)
}
