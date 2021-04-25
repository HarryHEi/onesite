package http

import (
	"fmt"
	"net"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"

	"onesite/core/config"
	"onesite/core/dao"
	"onesite/core/server/http/service"
	"onesite/core/worker"
)

type ServerConfig struct {
	Bind string `toml:"bind"`
	Port int    `toml:"port"`
	Rate int64  `toml:"rate"`
}

type Server struct {
	Svc *service.Service
	Cfg ServerConfig
	S   *gin.Engine
	M   *melody.Melody
	Dao *dao.Dao
	w   *worker.Worker
}

func NewHttpServer() (*Server, error) {
	var cfg ServerConfig
	_, err := toml.DecodeFile(config.GetCfgPath("server.toml"), &cfg)
	if err != nil {
		return nil, err
	}

	d, err := dao.NewDao()
	if err != nil {
		return nil, nil
	}

	w := worker.NewWorker(d)

	svc := service.NewService(d, w)

	s := &Server{
		svc,
		cfg,
		gin.New(),
		melody.New(),
		d,
		w,
	}
	err = InitRouter(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) Run() error {
	lis, err := net.Listen(
		"tcp",
		fmt.Sprintf("%s:%d", s.Cfg.Bind, s.Cfg.Port),
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

func RunHttpServer() error {
	s, err := NewHttpServer()
	if err != nil {
		return err
	}
	return s.Run()
}
