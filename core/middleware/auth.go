package middleware

import (
	"errors"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"onesite/common/rest"
	"onesite/core/config"
	"onesite/core/dao"
	"onesite/core/model"
)

var (
	authMiddleware *jwt.GinJWTMiddleware
	identityKey    = "username"
)

type AuthPayload struct {
	Username string
}

type LoginForm struct {
	Username string
	Password string
}

func InitAuthMiddleware() (err error) {
	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "jwt",
		Key:         []byte(config.CoreCfg.SecretKey),
		Timeout:     config.CoreCfg.Auth.Timeout.Duration,
		MaxRefresh:  config.CoreCfg.Auth.Timeout.Duration,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*AuthPayload); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &AuthPayload{
				Username: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginForm LoginForm
			if err := c.ShouldBind(&loginForm); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginForm.Username
			password := loginForm.Password

			user, err := dao.Authorization(username, password)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return &AuthPayload{user.Username}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*AuthPayload); ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			rest.Unauthorized(c, message)
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	return err
}

func GetAuthMiddleware() *jwt.GinJWTMiddleware {
	if authMiddleware == nil {
		panic("call InitAuthMiddleware before GetAuthMiddleware")
	}

	return authMiddleware
}

// 解析用户实例，保存到key="userInstance"
func ParseUserMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		_, exist := c.Get("userInstance")
		if exist {
			return
		}

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
		c.Set("userInstance", user)
	}
}

// 从上下文获取设置的用户实例
func ParseUser(c *gin.Context) *model.User {
	user, exist := c.Get("userInstance")
	if !exist {
		return nil
	}
	userInstance, ok := user.(*model.User)
	if !ok {
		return nil
	}
	return userInstance
}
