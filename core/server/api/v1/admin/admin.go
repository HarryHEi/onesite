package admin

import (
	"github.com/gin-gonic/gin"

	"onesite/common/rest"
	"onesite/core/dao"
	"onesite/core/model"
)

func ListUsers() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request rest.PaginationQueryParams
		err := c.ShouldBind(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}

		count, users, err := dao.ListUser(
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

func CreateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request CreateUserRequest
		err := c.ShouldBind(&request)
		if err != nil {
			rest.BadRequest(c, err)
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
		var request rest.PKDetailUri
		err := c.ShouldBindUri(&request)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		err = dao.DeleteUser(request.PK)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.NoContent(c)
	}
}

func PatchUpdateUser() func(c *gin.Context) {
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
		err = dao.UpdateUser(pkRequest.PK, request.Fields())
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.NoContent(c)
	}
}
