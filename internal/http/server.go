package HttpServer

import (
	"fmt"
	"github.com/aerosystems/customer-service/internal/infrastructure/rest"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const webPort = 80

type Server struct {
	log             *logrus.Logger
	echo            *echo.Echo
	customerHandler *rest.CustomerHandler
	tokenService    TokenService
}

func NewServer(
	log *logrus.Logger,
	customerHandler *rest.CustomerHandler,
	tokenService TokenService,

) *Server {
	return &Server{
		log:             log,
		echo:            echo.New(),
		customerHandler: customerHandler,
		tokenService:    tokenService,
	}
}

func (s *Server) Run() error {
	s.setupMiddleware()
	s.setupRoutes()
	s.log.Infof("starting HTTP server customer-service on port %d\n", webPort)
	return s.echo.Start(fmt.Sprintf(":%d", webPort))
}
