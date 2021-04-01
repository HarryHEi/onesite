package admin

import "onesite/core/model"

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
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"is_admin"`
}
