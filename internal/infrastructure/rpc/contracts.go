package RPCServer

import "github.com/aerosystems/customer-service/internal/models"

type CustomerUsecase interface {
	CreateUser() (*models.Customer, error)
}
