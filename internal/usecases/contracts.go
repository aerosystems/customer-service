package usecases

import (
	"context"
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	GetByUuid(ctx context.Context, uuid uuid.UUID) (*models.Customer, error)
	Create(ctx context.Context, user *models.Customer) error
	Update(ctx context.Context, user *models.Customer) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}

type ProjectRepository interface {
	CreateDefaultProject(customer *models.Customer) error
}

type SubsRepository interface {
	CreateFreeTrial(customer *models.Customer) error
	DeleteSubscription(customer *models.Customer) error
}
