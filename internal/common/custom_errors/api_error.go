package CustomErrors

import (
	"google.golang.org/grpc/codes"
	"net/http"
)

type ApiError struct {
	Message  string
	HttpCode int
	GrpcCode codes.Code
}

func (e ApiError) Error() string {
	return e.Message
}

var (
	ErrCustomerAlreadyExists = ApiError{"Customer already exists", http.StatusConflict, codes.AlreadyExists}
	ErrCustomerNotFound      = ApiError{"Customer not found", http.StatusNotFound, codes.NotFound}
)
