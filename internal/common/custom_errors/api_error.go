package CustomErrors

type APIError struct {
	Message  string
	HttpCode int
	RpcCode  int
}

func (e APIError) Error() string {
	return e.Message
}

var (
	ErrCustomerAlreadyExists = APIError{"Customer already exists", 409, 6}
	ErrCustomerNotFound      = APIError{"Customer not found", 404, 7}
)
