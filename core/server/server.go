package server

import "onesite/core/server/http"

func RunServer() (err error) {
	return http.RunHttpServer()
}
