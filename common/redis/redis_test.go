package redis_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"onesite/common/redis"
)

func TestInitRedis(t *testing.T) {
	_, err := redis.GetRedis()
	require.NotNil(t, err)
	err = redis.InitRedis()
	require.Nil(t, err)
	cli, _ := redis.GetRedis()
	require.NotNil(t, cli)
}
