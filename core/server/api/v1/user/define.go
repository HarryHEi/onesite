package user

import "onesite/core/model"

type InfoResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"is_admin"`
}

func InfoResponseFromUserModel(user *model.User) *InfoResponse {
	return &InfoResponse{
		Id:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		IsAdmin:  user.IsAdmin,
	}
}
