package CustomErrors

import (
	"net/http"
)

type APIError struct {
	Message  string
	HttpCode int
}

func (e APIError) Error() string {
	return e.Message
}

var (
	ErrCustomerAlreadyExists = APIError{"Customer already exists", http.StatusConflict}
	ErrCustomerNotFound      = APIError{"Customer not found", http.StatusNotFound}
)
