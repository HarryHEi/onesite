package sync

import (
	"onesite/core/config"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"onesite/core/dao/db"
)

func TestRDLock_Lock(t *testing.T) {
	config.CfgRootPath = "../../../configs"

	rd, err := db.NewRedis()
	require.Nil(t, err)

	key := "test_key"
	value := uuid.New().String()

	l := NewRDLock(rd.Db)

	locked := l.Lock(key, value, 10000)
	require.Equal(t, locked, true)

	relocked := l.Lock(key, value, 10000)
	require.Equal(t, relocked, false)

	wrongUnlock := l.UnLock(key, "wrong value")
	require.Equal(t, wrongUnlock, false)

	unlocked := l.UnLock(key, value)
	require.Equal(t, unlocked, true)
}
