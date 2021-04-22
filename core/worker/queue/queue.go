package queue

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Queue interface {
	GetTopic(topic string) interface{}
	PutTopic(topic string, v interface{})
}

type RedisQueue struct {
	rd *redis.Client
}

func NewRedisQueue(rd *redis.Client) *RedisQueue {
	return &RedisQueue{rd}
}

func (r *RedisQueue) GetTopic(topic string) interface{} {
	cmd := r.rd.BLPop(context.Background(), 0, topic)
	return cmd.Val()[1]
}

func (r *RedisQueue) PutTopic(topic string, v interface{}) {
	r.rd.RPush(context.Background(), topic, v)
}

func (r *RedisQueue) ClearTopic(topic string) {
	r.rd.Del(context.Background(), topic)
}
