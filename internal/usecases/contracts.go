package usecases

import (
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	GetByUuid(uuid uuid.UUID) (*models.Customer, error)
	Create(user *models.Customer) error
	Update(user *models.Customer) error
	Delete(user *models.Customer) error
}

type ProjectRepository interface {
	CreateDefaultProject(customer *models.Customer) error
}

type SubsRepository interface {
	CreateFreeTrial(customer *models.Customer) error
	DeleteSubscription(customer *models.Customer) error
}
