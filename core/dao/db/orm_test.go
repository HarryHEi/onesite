package db

import (
	"testing"

	"github.com/stretchr/testify/require"

	"onesite/core/config"
)

func TestInitOrm(t *testing.T) {
	config.CfgRootPath = "../../../configs"

	_, err := NewOrm()
	require.Nil(t, err)
}
