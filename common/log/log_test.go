package log_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"onesite/common/log"
)

func TestInitLogger(t *testing.T) {
	err := log.InitLogger()
	require.Nil(t, err)
}
