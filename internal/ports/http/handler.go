package HTTPServer

type Handler struct {
	customerUsecase CustomerUsecase
}

func NewHandler(
	customerUsecase CustomerUsecase,
) *Handler {
	return &Handler{
		customerUsecase: customerUsecase,
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}
