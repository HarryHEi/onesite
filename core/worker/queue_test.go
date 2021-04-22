package worker_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"onesite/common/redis"
	"onesite/core/worker"
)

func TestRedisQueue(t *testing.T) {
	_ = redis.InitRedis()
	rd, _ := redis.GetRedis()
	q := worker.NewRedisQueue(rd, "my-topic-1")
	q.Clear()
	for number := 1; number < 10; number++ {
		q.Put(number)
	}
	for number := 1; number < 10; number++ {
		v := q.Get()
		v, ok := v.(string)
		require.Equal(t, ok, true)
		require.Equal(t, v, fmt.Sprintf("%d", number))
	}
}
