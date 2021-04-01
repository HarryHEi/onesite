package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"onesite/common/rest"
	"onesite/core/dao"
	"onesite/core/model"
)

func ListUsers() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PaginationQueryParams
		err := c.Bind(&request)
		if err != nil {
			rest.BadRequest(c, errors.New("invalid params"))
			return
		}

		count, users, err := dao.ListUser([]string{"id", "username", "name", "is_admin"}, 1, 10)
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

func CreateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request CreateUserRequest
		err := c.Bind(&request)
		if err != nil {
			rest.BadRequest(c, errors.New("invalid params"))
			return
		}

		user, err := dao.CreateUser(&model.User{
			Username: request.Username,
			Password: dao.GeneratePassword(request.Password),
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

func DeleteUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		pk := c.Param("pk")
		err := dao.DeleteUser(pk)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.NoContent(c)
	}
}

func PatchUpdateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		pk := c.Param("pk")
		request := make(map[string]interface{})
		err := c.Bind(&request)
		if err != nil {
			rest.BadRequest(c, errors.New("invalid params"))
			return
		}

		err = dao.UpdateUser(pk, request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.NoContent(c)
	}
}
