package auth

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"onesite/common/log"
	"onesite/common/rest"
	"onesite/core/dao"
)

type LoginForm struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func Login() func(c *gin.Context) {
	return func(c *gin.Context) {
		var loginForm LoginForm
		err := c.Bind(&loginForm)
		if err != nil {
			log.Error("Bind login form failed", zap.Error(err))
			rest.BadRequest(c, err)
			return
		}

		user, err := dao.Authorization(loginForm.Username, loginForm.Password)
		if err != nil {
			log.Error("Login failed", zap.Error(err))
			rest.BadRequest(c, err)
			return
		}

		rest.Success(c, user)
	}
}
