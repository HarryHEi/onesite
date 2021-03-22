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
