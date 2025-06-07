package entities

import (
	"net/http"

	"google.golang.org/grpc/codes"

	"github.com/aerosystems/common-service/customerrors"
)

var (
	ErrCustomerAlreadyExists = customerrors.InternalError{Message: "Customer already exists", HttpCode: http.StatusConflict, GrpcCode: codes.AlreadyExists}
	ErrCustomerNotFound      = customerrors.InternalError{Message: "Customer not found", HttpCode: http.StatusNotFound, GrpcCode: codes.NotFound}
)
