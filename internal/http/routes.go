package HttpServer

import "github.com/aerosystems/customer-service/internal/models"

func (s *Server) setupRoutes() {
	s.echo.GET("/v1/customers", s.customerHandler.GetCustomer, s.AuthTokenMiddleware(models.CustomerRole))
}
