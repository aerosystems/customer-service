package app

import (
	"github.com/sirupsen/logrus"

	HTTPServer "github.com/aerosystems/customer-service/internal/ports/http"
)

type Server struct {
	log  *logrus.Logger
	cfg  *Config
	http *HTTPServer.Server
}

func NewServer(
	log *logrus.Logger,
	cfg *Config,
	httpServer *HTTPServer.Server,
) *Server {
	return &Server{
		log:  log,
		cfg:  cfg,
		http: httpServer,
	}
}
