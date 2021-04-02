package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"onesite/common/orm"
	"onesite/core/model"
)

type AuthTestSuite struct {
	suite.Suite
	OrmDb *gorm.DB
}

func (s *AuthTestSuite) SetupTest() {
	err := orm.InitOrm(
		orm.DriverName("mysql"),
		orm.Dsn("herui:Admin@123@tcp(172.172.177.191:3306)/onesite_dev?charset=utf8mb4&parseTime=true"),
	)
	require.Nil(s.T(), err)
	s.OrmDb, err = orm.GetDb()
	require.Nil(s.T(), err)
	s.OrmDb.Raw("t")
	err = s.OrmDb.AutoMigrate(&model.User{})
	require.Nil(s.T(), err)
}

func (s *AuthTestSuite) TestCRUD() {
	// Create
	createUser := model.User{
		Username: "username",
		Name:     "name",
		IsAdmin:  true,
	}
	s.OrmDb.Create(&createUser)
	require.NotEmpty(s.T(), createUser.ID)

	// Query
	queryUser := model.User{}
	s.OrmDb.Model(&model.User{}).First(&queryUser, createUser.ID)
	require.Equal(s.T(), createUser.Name, queryUser.Name)

	// Update And Query
	s.OrmDb.Model(&queryUser).Update("Name", "张三")
	queryAfterUpdateUser := model.User{}
	s.OrmDb.Model(&model.User{}).First(&queryAfterUpdateUser, createUser.ID)
	require.Equal(s.T(), queryAfterUpdateUser.Name, "张三")

	// Delete
	s.OrmDb.Model(&model.User{}).Unscoped().Delete(&model.User{}, createUser.ID)
	queryAfterDeleteUser := model.User{}
	result := s.OrmDb.Model(&model.User{}).Find(queryAfterDeleteUser, createUser.ID)
	require.Empty(s.T(), result.RowsAffected)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
