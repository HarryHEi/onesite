package redis_test

import (
	"testing"

	"github.com/stretchr/testify/require"

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
