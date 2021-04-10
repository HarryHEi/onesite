package dao

import (
	"crypto/md5"
	"errors"
	"fmt"

	"onesite/core/model"
)

// 生成MD5加密的密码
func GeneratePassword(password string) string {
	pass := []byte(password)
	return fmt.Sprintf("%x", md5.Sum(pass))
}

// 验证比较密码明文和MD5加密的密文
func CheckPassword(password, genPass string) bool {
	return GeneratePassword(password) == genPass
}

// 通过用户名查询用户信息
func QueryUser(username string) (*model.User, error) {
	daoIns, err := GetDao()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	ret := daoIns.Db.Model(&model.User{}).Where("username=?", username).First(user)
	if ret.Error != nil {
		return nil, fmt.Errorf("query user failed %e", ret.Error)
	}
	return user, nil
}

// 通过用户名密码验证用户
func Authorization(username, password string) (*model.User, error) {
	user, err := QueryUser(username)
	if err != nil {
		return nil, err
	}
	if user.Username != username {
		return nil, fmt.Errorf("unknown user %s", username)
	}
	if !CheckPassword(password, user.Password) {
		return nil, errors.New("wrong password")
	}
	return user, nil
}

// 分页查询用户
func ListUser(fields []string, page, pageSize int) (count int64, users []model.User, err error) {
	daoIns, err := GetDao()
	if err != nil {
		return 0, nil, err
	}

	if page <= 0 || pageSize <= 0 {
		page = 1
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	ret := daoIns.Db.Model(&model.User{}).Select(fields).Count(&count)
	if ret.Error != nil {
		return 0, nil, ret.Error
	}
	if count == 0 {
		return 0, nil, nil
	}

	ret = daoIns.Db.Model(&model.User{}).Select(fields).Offset(offset).Limit(pageSize).Find(&users)
	if ret.Error != nil {
		return 0, nil, ret.Error
	}
	return count, users, nil
}

// 新增用户
func CreateUser(user *model.User) (*model.User, error) {
	daoIns, err := GetDao()
	if err != nil {
		return nil, err
	}

	ret := daoIns.Db.Create(user)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return user, nil
}

// 删除用户
func DeleteUser(pk interface{}) error {
	daoIns, err := GetDao()
	if err != nil {
		return err
	}

	ret := daoIns.Db.Model(&model.User{}).Unscoped().Delete(model.User{}, pk)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

// 更新用户
func UpdateUser(pk, v interface{}) error {
	daoIns, err := GetDao()
	if err != nil {
		return err
	}

	ret := daoIns.Db.Model(&model.User{}).Where("id = ?", pk).Updates(v)
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

// 创建管理员账户
func CreateSuperuserIfNotExists(username, password string) error {
	_, err := CreateUser(&model.User{
		Username: username,
		Password: GeneratePassword(password),
		Name:     "管理员",
		IsAdmin:  true,
	})
	return err
}
