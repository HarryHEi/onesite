package sync_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"onesite/common/redis"
	"onesite/common/sync"
)

func TestRDLock_Lock(t *testing.T) {
	err := redis.InitRedis()
	require.Nil(t, err)

	r, err := redis.GetRedis()
	require.NotNil(t, r)
	require.Nil(t, err)

	key := "test_key"
	value := uuid.New().String()

	sync.InitRDLock(r)
	l, _ := sync.GetRDLock()

	locked := l.Lock(key, value, 10000)
	require.Equal(t, locked, true)

	relocked := l.Lock(key, value, 10000)
	require.Equal(t, relocked, false)

	wrongUnlock := l.UnLock(key, "wrong value")
	require.Equal(t, wrongUnlock, false)

	unlocked := l.UnLock(key, value)
	require.Equal(t, unlocked, true)
}
