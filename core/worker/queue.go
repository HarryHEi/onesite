package worker

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Queue interface {
	Get() interface{}
	Put(interface{})
}

type RedisQueue struct {
	Topic string

	rd *redis.Client
}

func NewRedisQueue(rd *redis.Client, topic string) *RedisQueue {
	return &RedisQueue{
		topic,
		rd,
	}
}

func (r *RedisQueue) Get() interface{} {
	cmd := r.rd.BLPop(context.Background(), 0, r.Topic)
	return cmd.Val()[1]
}

func (r *RedisQueue) Put(v interface{}) {
	r.rd.RPush(context.Background(), r.Topic, v)
}

func (r *RedisQueue) Clear() {
	r.rd.Del(context.Background(), r.Topic)
}
