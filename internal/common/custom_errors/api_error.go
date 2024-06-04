package CustomErrors

import (
	"net/http"
)

type ApiError struct {
	Message  string
	HttpCode int
}

func (e ApiError) Error() string {
	return e.Message
}

var (
	ErrCustomerAlreadyExists = ApiError{"Customer already exists", http.StatusConflict}
	ErrCustomerNotFound      = ApiError{"Customer not found", http.StatusNotFound}
)
