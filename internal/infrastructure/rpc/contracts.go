package RpcServer

import "github.com/aerosystems/customer-service/internal/models"

type CustomerUsecase interface {
	CreateCustomer() (*models.Customer, error)
}
