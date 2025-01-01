package usecases

import (
	"context"
	"github.com/aerosystems/customer-service/internal/domain"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Customer, error)
	Create(ctx context.Context, user *domain.Customer) error
	Update(ctx context.Context, user *domain.Customer) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}

type SubscriptionEventsAdapter interface {
	PublishCreateFreeTrialEvent(customerUuid uuid.UUID) error
}
