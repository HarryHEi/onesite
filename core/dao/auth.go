package dao

import (
	"crypto/md5"
	"errors"
	"fmt"

	"onesite/core/model"
)

// GeneratePassword 生成MD5加密的密码
func (dao *Dao) GeneratePassword(password string) string {
	pass := []byte(password)
	return fmt.Sprintf("%x", md5.Sum(pass))
}

// CheckPassword 验证比较密码明文和MD5加密的密文
func (dao *Dao) CheckPassword(password, genPass string) bool {
	return dao.GeneratePassword(password) == genPass
}

// QueryUser 通过用户名查询用户信息
func (dao *Dao) QueryUser(username string) (*model.User, error) {
	user := &model.User{}
	ret := dao.Orm.Db.Model(&model.User{}).Where("username=?", username).First(user)
	if ret.Error != nil {
		return nil, fmt.Errorf("query user failed %e", ret.Error)
	}
	return user, nil
}

// Authorization 通过用户名密码验证用户
func (dao *Dao) Authorization(username, password string) (*model.User, error) {
	user, err := dao.QueryUser(username)
	if err != nil {
		return nil, err
	}
	if user.Username != username {
		return nil, fmt.Errorf("unknown user %s", username)
	}
	if !dao.CheckPassword(password, user.Password) {
		return nil, errors.New("wrong password")
	}
	return user, nil
}

// ListUser 分页查询用户
func (dao *Dao) ListUser(fields []string, page, pageSize int) (count int64, users []model.User, err error) {
	if page <= 0 || pageSize <= 0 {
		page = 1
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	ret := dao.Orm.Db.Model(&model.User{}).Select(fields).Count(&count)
	if ret.Error != nil {
		return 0, nil, ret.Error
	}
	if count == 0 {
		return 0, nil, nil
	}

	ret = dao.Orm.Db.Model(&model.User{}).Select(fields).Offset(offset).Limit(pageSize).Find(&users)
	if ret.Error != nil {
		return 0, nil, ret.Error
	}
	return count, users, nil
}

// CreateUser 新增用户
func (dao *Dao) CreateUser(user *model.User) (*model.User, error) {
	ret := dao.Orm.Db.Create(user)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return user, nil
}

// DeleteUser 删除用户
func (dao *Dao) DeleteUser(pk interface{}) error {
	ret := dao.Orm.Db.Model(&model.User{}).Unscoped().Delete(model.User{}, pk)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

// UpdateUser 更新用户
func (dao *Dao) UpdateUser(pk, v interface{}) error {
	ret := dao.Orm.Db.Model(&model.User{}).Where("id = ?", pk).Updates(v)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

// CreateSuperuserIfNotExists 创建管理员账户
func (dao *Dao) CreateSuperuserIfNotExists(username, password string) error {
	_, err := dao.CreateUser(&model.User{
		Username: username,
		Password: dao.GeneratePassword(password),
		Name:     "管理员",
		IsAdmin:  true,
	})
	return err
}

// QueryUserById 通过用户名查询用户信息
func (dao *Dao) QueryUserById(pk interface{}) (*model.User, error) {
	user := &model.User{}
	ret := dao.Orm.Db.Model(&model.User{}).First(user, pk)
	if ret.Error != nil {
		return nil, fmt.Errorf("query user failed %e", ret.Error)
	}
	return user, nil
}
