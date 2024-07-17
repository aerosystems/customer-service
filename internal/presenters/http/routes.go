package HttpServer

func (s *Server) setupRoutes() {
	s.echo.POST("/v1/customers", s.customerHandler.CreateCustomer)
}
