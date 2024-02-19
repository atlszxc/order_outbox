package server

import (
	"github.com/gin-gonic/gin"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PATCH  = "PATCH"
)

type Server struct {
	S      *gin.Engine
	Config Config
}

func (s *Server) AddRoute(method, path string, handler gin.HandlerFunc) {
	switch method {
	case GET:
		s.S.GET(path, handler)
	case POST:
		s.S.POST(path, handler)
	case DELETE:
		s.S.DELETE(path, handler)
	case PATCH:
		s.S.PATCH(path, handler)
	}
}

func (s *Server) Listen() error {
	err := s.S.Run(s.Config.Port)
	return err
}

func InitServer(port, mode string) *Server {
	gin.SetMode(mode)
	return &Server{
		S: gin.Default(),
		Config: Config{
			Port: port,
			Mode: mode,
		},
	}
}
