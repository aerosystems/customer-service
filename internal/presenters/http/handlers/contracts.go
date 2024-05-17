package handlers

import (
	"github.com/aerosystems/customer-service/internal/models"
)

type CustomerUsecase interface {
	CreateCustomer(uuid string) (*models.Customer, error)
}
