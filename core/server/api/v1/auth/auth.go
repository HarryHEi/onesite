package auth

import (
	"errors"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"onesite/common/rest"
	"onesite/core/dao"
)

type LoginForm struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func UserInfo() func(c *gin.Context) {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		username, ok := claims["username"].(string)
		if !ok {
			rest.BadRequest(c, errors.New("invalid username"))
		}
		user, err := dao.QueryUser(username)
		if err != nil {
			rest.BadRequest(c, err)
			return
		}
		rest.Success(c, user)
	}
}
