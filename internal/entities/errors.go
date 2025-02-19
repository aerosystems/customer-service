package entities

import (
	"github.com/aerosystems/common-service/customerrors"
	"google.golang.org/grpc/codes"
	"net/http"
)

var (
	ErrCustomerAlreadyExists = customerrors.InternalError{Message: "Customer already exists", HttpCode: http.StatusConflict, GrpcCode: codes.AlreadyExists}
	ErrCustomerNotFound      = customerrors.InternalError{Message: "Customer not found", HttpCode: http.StatusNotFound, GrpcCode: codes.NotFound}
)
