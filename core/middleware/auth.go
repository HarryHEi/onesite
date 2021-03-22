package middleware

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"onesite/common/rest"
	"onesite/core/config"
	"onesite/core/dao"
)

var (
	authMiddleware *jwt.GinJWTMiddleware
	identityKey    = "id"
)

type AuthPayload struct {
	Id uint
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
					identityKey: v.Id,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &AuthPayload{
				Id: claims[identityKey].(uint),
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

			return &AuthPayload{user.ID}, nil
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
