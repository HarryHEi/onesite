package service

import (
	"github.com/gin-gonic/gin"

	"onesite/core/dao"
	"onesite/core/model"
	"onesite/core/server/http/rest"
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

func (u *UpdateUserRequest) Fields(d *dao.Dao) map[string]interface{} {
	return map[string]interface{}{
		"username": u.Username,
		"name":     u.Name,
		"password": d.GeneratePassword(u.Password),
		"is_admin": u.IsAdmin,
	}
}

func (s *Service) ListUsers() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PaginationQueryParams
		err := c.ShouldBind(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		count, users, err := s.Dao.ListUser(
			[]string{"id", "username", "name", "is_admin"},
			request.Page,
			request.PageSize,
		)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.Success(c, rest.PaginationResponse{
			Count: count,
			Data:  UsersResponseFromUserModels(users),
		})
	}
}

func (s *Service) CreateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request CreateUserRequest
		err := c.ShouldBind(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		user, err := s.Dao.CreateUser(&model.User{
			Username: request.Username,
			Password: s.Dao.GeneratePassword(request.Password),
			Name:     request.Name,
			IsAdmin:  request.IsAdmin,
		})
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.Created(c, UserResponseFromUserModel(user))
	}
}

func (s *Service) DeleteUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PKDetailUri
		err := c.ShouldBindUri(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		err = s.Dao.DeleteUser(request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.NoContent(c)
	}
}

func (s *Service) PatchUpdateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var pkRequest rest.PKDetailUri
		err := c.ShouldBindUri(&pkRequest)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		var request UpdateUserRequest
		err = c.ShouldBind(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		err = s.Dao.UpdateUser(pkRequest.PK, request.Fields(s.Dao))
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.NoContent(c)
	}
}
