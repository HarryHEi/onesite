package sync

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RDLock Redis Distribute Lock
type RDLock struct {
	_redis *redis.Client
	ctx    context.Context
}

var (
	lockScript   = `return redis.call('SET', KEYS[1], ARGV[1], 'NX', 'PX', ARGV[2])`
	unlockScript = `
		if redis.call("get",KEYS[1]) == ARGV[1] then
		    return redis.call("del",KEYS[1])
		else
		    return 0
		end
	`
)

func NewRDLock(r *redis.Client) *RDLock {
	return &RDLock{
		r,
		context.Background(),
	}
}

func (r *RDLock) Lock(key, value string, timeoutMs int) bool {
	script := redis.NewScript(lockScript)
	cmd := script.Run(r.ctx, r._redis, []string{key}, value, timeoutMs)
	return cmd.Val() == "OK"
}

func (r *RDLock) UnLock(key, value string) bool {
	script := redis.NewScript(unlockScript)
	cmd := script.Run(r.ctx, r._redis, []string{key}, value)
	return cmd.Val().(int64) != 0
}

func (r *RDLock) TryLockUntil(key, value string, timeoutMs int, until time.Duration) bool {
	for {
		select {
		case <-time.After(until):
			return false
		case <-time.Tick(time.Second):
			if r.Lock(key, value, timeoutMs) {
				return true
			}
		}
	}
}
