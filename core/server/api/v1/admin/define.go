package admin

import (
	"onesite/core/dao"
	"onesite/core/model"
)

type UsersResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"is_admin"`
}

func UserResponseFromUserModel(user *model.User) *UsersResponse {
	return &UsersResponse{
		Id:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		IsAdmin:  user.IsAdmin,
	}
}

func UsersResponseFromUserModels(users []model.User) []*UsersResponse {
	usersResponse := make([]*UsersResponse, 0, len(users))
	for index := range users {
		usersResponse = append(usersResponse, UserResponseFromUserModel(&users[index]))
	}
	return usersResponse
}

type CreateUserRequest struct {
	Username string `json:"username" form:"username" binding:"required,gte=3,lte=32"`
	Password string `json:"password" form:"password" binding:"required,gte=6,lte=32"`
	Name     string `json:"name" form:"name" binding:"required,gte=1,lte=64"`
	IsAdmin  bool   `json:"is_admin" form:"is_admin"`
}

type UpdateUserRequest struct {
	Username string `json:"username" form:"username" binding:"required,gte=3,lte=32"`
	Name     string `json:"name" form:"name" binding:"required,gte=1,lte=64"`
	Password string `json:"password" form:"password" binding:"required,gte=6,lte=32"`
	IsAdmin  bool   `json:"is_admin" form:"is_admin"`
}

func (u *UpdateUserRequest) Fields() map[string]interface{} {
	return map[string]interface{}{
		"username": u.Username,
		"name":     u.Name,
		"password": dao.GeneratePassword(u.Password),
		"is_admin": u.IsAdmin,
	}
}
