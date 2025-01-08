package HTTPServer

func (s *Server) setupRoutes() {
	s.echo.POST("/v1/firebase/create-customer", s.firebaseHandler.CreateCustomer)
}
