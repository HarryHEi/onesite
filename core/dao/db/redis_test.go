package db

import (
	"testing"

	"github.com/stretchr/testify/require"

	"onesite/core/config"
)

func TestInitRedis(t *testing.T) {
	config.CfgRootPath = "../../../configs"

	_, err := NewRedis()
	require.Nil(t, err)
}
