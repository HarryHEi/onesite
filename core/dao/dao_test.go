package dao

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"onesite/core/config"
	"onesite/core/model"
)

type AuthTestSuite struct {
	suite.Suite
	daoIns   *Dao
	username string
	password string
}

// 初始化Dao，建立一个测试用户。
func (s *AuthTestSuite) SetupTest() {
	config.CfgRootPath = "../../configs"

	dao, err := NewDao()
	require.Nil(s.T(), err)

	s.daoIns = dao
	require.NotNil(s.T(), s.daoIns)

	s.username = "20210227"
	s.password = "my pass"

	// Create User
	createUser := model.User{
		Username: s.username,
		Name:     "测试用户",
		Password: s.daoIns.GeneratePassword(s.password),
		IsAdmin:  true,
	}
	s.daoIns.Orm.Db.Create(&createUser)
	require.NotEmpty(s.T(), createUser.ID)
}

// 测试密码加密
func (s *AuthTestSuite) TestGeneratePassword() {
	s1 := s.daoIns.GeneratePassword("test")
	s2 := s.daoIns.GeneratePassword("test")
	s3 := s.daoIns.GeneratePassword("test2")
	require.Equal(s.T(), s1, s2)
	require.NotEqual(s.T(), s1, s3)
}

// 测试用户认证
func (s *AuthTestSuite) TestAuthorization() {
	authUser, err := s.daoIns.Authorization(s.username, s.password)
	require.NotNil(s.T(), authUser)
	require.Nil(s.T(), err)
	require.Equal(s.T(), authUser.Username, s.username)
}

// 删除测试用户
func (s *AuthTestSuite) TearDownTest() {
	// Delete User
	s.daoIns.Orm.Db.Model(&model.User{}).Unscoped().Where("username=?", s.username).Delete(&model.User{})
	queryAfterDeleteUser := model.User{}
	result := s.daoIns.Orm.Db.Model(&model.User{}).Where("username=?", s.username).Find(queryAfterDeleteUser)
	require.Empty(s.T(), result.RowsAffected)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
