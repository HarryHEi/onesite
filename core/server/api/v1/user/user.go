package user

import (
	"github.com/gin-gonic/gin"

	"onesite/common/rest"
	"onesite/core/middleware"
)

func Info() func(c *gin.Context) {
	return func(c *gin.Context) {
		user := middleware.ParseUser(c)
		rest.Success(c, InfoResponseFromUserModel(user))
	}
}
