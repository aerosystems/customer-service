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

type SubscriptionEventsAdapter interface {
	PublishCreateFreeTrialEvent(customerUuid uuid.UUID) error
}
