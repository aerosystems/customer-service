package usecases

import (
	"context"

	"github.com/google/uuid"

	"github.com/aerosystems/customer-service/internal/entities"
)

type CustomerRepository interface {
	GetByCustomerUUID(ctx context.Context, customerUUID uuid.UUID) (*entities.Customer, error)
	GetByFirebaseUID(ctx context.Context, firebaseUID string) (*entities.Customer, error)
	Create(ctx context.Context, customer *entities.Customer) error
	Upsert(ctx context.Context, customer *entities.Customer) error
	Delete(ctx context.Context, customerUUID uuid.UUID) error
}

type SubscriptionAdapter interface {
	CreateFreeTrialSubscription(ctx context.Context, customerUUID uuid.UUID) (*SubscriptionDTO, error)
	DeleteSubscription(ctx context.Context, subscriptionUUID uuid.UUID) error
}

type ProjectAdapter interface {
	CreateDefaultProject(ctx context.Context, customerUUID uuid.UUID) (uuid.UUID, string, error)
	DeleteProject(ctx context.Context, projectUUID uuid.UUID) error
}

type CheckmailAdapter interface {
	CreateAccess(ctx context.Context, projectToken string, subscriptionDTO *SubscriptionDTO) error
}

type FirebaseAuthAdapter interface {
	SetCustomUserClaims(ctx context.Context, uid string, customerUUID uuid.UUID) error
}
