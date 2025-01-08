package usecases

import (
	"context"
	"github.com/aerosystems/customer-service/internal/domain"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	GetByCustomerUUID(ctx context.Context, customerUUID uuid.UUID) (*domain.Customer, error)
	GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.Customer, error)
	Create(ctx context.Context, customer *domain.Customer) error
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, customerUUID uuid.UUID) error
}

type SubscriptionAdapter interface {
	CreateFreeTrialSubscription(ctx context.Context, customerUUID uuid.UUID) (subscriptionUUID uuid.UUID, err error)
	DeleteSubscription(ctx context.Context, subscriptionUUID uuid.UUID) error
}

type ProjectAdapter interface {
	CreateDefaultProject(ctx context.Context, customerUUID uuid.UUID) (projectUUID uuid.UUID, err error)
	DeleteProject(ctx context.Context, projectUUID uuid.UUID) error
}

type FirebaseAuthAdapter interface {
	SetClaimCustomerUUID(ctx context.Context, uid string, customerUUID uuid.UUID) error
}
