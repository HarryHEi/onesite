package http

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"onesite/core/config"
	"onesite/core/server/router"
)

func NewHttpServer() *http.Server {
	_s := gin.New()

	router.InitRouter(_s)

	return &http.Server{
		Handler: _s,
	}
}

func RunHttpServer() (err error) {
	lis, err := net.Listen(
		"tcp",
		fmt.Sprintf("%s:%d", config.CoreCfg.Server.Bind, config.CoreCfg.Server.Port),
	)
	if err != nil {
		return err
	}
	httpServer := NewHttpServer()
	err = httpServer.Serve(lis)
	return err
}
