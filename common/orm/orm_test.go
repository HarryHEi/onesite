package orm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"onesite/common/orm"
)

func TestInitOrm(t *testing.T) {
	_, err := orm.GetDb()
	require.NotNil(t, err)
	err = orm.InitOrm(
		orm.DriverName("mysql"),
		orm.Dsn("herui:Admin@123@tcp(172.172.177.191:3306)/onesite_dev?charset=utf8mb4&parseTime=true"),
		orm.MaxOpenConn(6),
		orm.MaxIdleConn(3),
	)
	require.Nil(t, err)
	_, err = orm.GetDb()
	require.Nil(t, err)
}
