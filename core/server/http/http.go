package http

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"

	"onesite/core/config"
)

type Service struct {
	S *gin.Engine
	M *melody.Melody
}

func NewHttpService() *Service {
	_s := &Service{
		gin.New(),
		melody.New(),
	}

	InitRouter(_s)

	return _s
}

func (s *Service) Run() error {
	lis, err := net.Listen(
		"tcp",
		fmt.Sprintf("%s:%d", config.CoreCfg.Server.Bind, config.CoreCfg.Server.Port),
	)
	if err != nil {
		return err
	}
	server := http.Server{
		Handler: s.S,
	}
	err = server.Serve(lis)
	return err
}

func RunHttpServer() (err error) {
	service := NewHttpService()
	err = service.Run()
	return err
}
