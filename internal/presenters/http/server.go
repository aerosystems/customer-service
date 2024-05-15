package HttpServer

import (
	"fmt"
	"github.com/aerosystems/customer-service/internal/presenters/http/handlers"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const webPort = 80

type Server struct {
	log             *logrus.Logger
	echo            *echo.Echo
	customerHandler *handlers.CustomerHandler
}

func NewServer(
	log *logrus.Logger,
	customerHandler *handlers.CustomerHandler,

) *Server {
	return &Server{
		log:             log,
		echo:            echo.New(),
		customerHandler: customerHandler,
	}
}

func (s *Server) Run() error {
	s.setupMiddleware()
	s.setupRoutes()
	s.log.Infof("starting HTTP server customer-service on port %d\n", webPort)
	return s.echo.Start(fmt.Sprintf(":%d", webPort))
}