package orm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"onesite/common/orm"
)

func TestInitOrm(t *testing.T) {
	_, err := orm.GetDb()
	require.NotNil(t, err)
	err = orm.InitOrm()
	require.Nil(t, err)
	_, err = orm.GetDb()
	require.Nil(t, err)
}
