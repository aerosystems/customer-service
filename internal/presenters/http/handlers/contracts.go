package handlers

import "github.com/aerosystems/customer-service/internal/models"

type CustomerUsecase interface {
	GetUserByUuid(uuid string) (*models.Customer, error)
}
