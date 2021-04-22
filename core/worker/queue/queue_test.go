package queue_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"onesite/common/redis"
	"onesite/core/worker/queue"
)

func TestRedisQueue(t *testing.T) {
	_ = redis.InitRedis()
	rd, _ := redis.GetRedis()
	const TOPIC = "my-topic-1"
	q := queue.NewRedisQueue(rd)
	q.ClearTopic(TOPIC)
	for number := 1; number < 10; number++ {
		q.PutTopic(TOPIC, number)
	}
	for number := 1; number < 10; number++ {
		v := q.GetTopic(TOPIC)
		v, ok := v.(string)
		require.Equal(t, ok, true)
		require.Equal(t, v, fmt.Sprintf("%d", number))
	}
}
