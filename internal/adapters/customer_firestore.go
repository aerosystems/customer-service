package adapters

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	CustomErrors "github.com/aerosystems/customer-service/internal/common/custom_errors"
	"github.com/aerosystems/customer-service/internal/domain"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const (
	customersCollectionName = "customers"
)

type FirestoreCustomerRepo struct {
	client *firestore.Client
}

func NewFirestoreCustomerRepo(client *firestore.Client) *FirestoreCustomerRepo {
	return &FirestoreCustomerRepo{
		client: client,
	}
}

type Customer struct {
	UUID        string     `firestore:"uuid"`
	Email       string     `firestore:"email"`
	FirebaseUID string     `firestore:"firebase_uid"`
	CreatedAt   time.Time  `firestore:"created_at"`
	DeleteAt    *time.Time `firestore:"delete_at,omitempty"`
}

func (c *Customer) ToModel() *domain.Customer {
	return &domain.Customer{
		UUID:        uuid.MustParse(c.UUID),
		Email:       c.Email,
		FirebaseUID: c.FirebaseUID,
		CreatedAt:   c.CreatedAt,
	}
}

func CustomerToFirestore(customer *domain.Customer) *Customer {
	return &Customer{
		UUID:        customer.UUID.String(),
		Email:       customer.Email,
		FirebaseUID: customer.FirebaseUID,
		CreatedAt:   customer.CreatedAt,
		DeleteAt:    nil,
	}
}

func (fcr *FirestoreCustomerRepo) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Customer, error) {
	docRef := fcr.client.Collection(customersCollectionName).Doc(uuid.String())
	doc, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, CustomErrors.ErrCustomerNotFound
		}
		return nil, fmt.Errorf("could not get customer from Firestore repository: %w", err)
	}

	var customer Customer
	if err = doc.DataTo(&customer); err != nil {
		return nil, fmt.Errorf("could not convert Firestore data to customer: %w", err)
	}

	if customer.DeleteAt != nil {
		return nil, CustomErrors.ErrCustomerNotFound
	}

	return customer.ToModel(), nil
}

func (fcr *FirestoreCustomerRepo) GetByFirebaseUID(_ context.Context, firebaseUID string) (*domain.Customer, error) {
	query := fcr.client.Collection(customersCollectionName).Where("firebase_uid", "==", firebaseUID).Limit(1)

	iter := query.Documents(context.Background())
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
			return nil, err
			break
		}
		customer := Customer{}
		err = doc.DataTo(&customer)
		if err != nil {
			return nil, err
			break
		}
		if customer.DeleteAt == nil {
			return customer.ToModel(), nil
			break
		}
	}
	return nil, CustomErrors.ErrCustomerNotFound
}

func (fcr *FirestoreCustomerRepo) Create(ctx context.Context, customer *domain.Customer) error {
	currentCustomer, err := fcr.GetByUUID(ctx, customer.UUID)
	if err != nil && !errors.Is(err, CustomErrors.ErrCustomerNotFound) {
		return err
	}
	if currentCustomer != nil {
		return CustomErrors.ErrCustomerAlreadyExists
	}

	_, err = fcr.client.Collection(customersCollectionName).Doc(customer.UUID.String()).Set(ctx, CustomerToFirestore(customer))
	return err
}

func (fcr *FirestoreCustomerRepo) Update(ctx context.Context, customer *domain.Customer) error {
	_, err := fcr.client.Collection(customersCollectionName).Doc(customer.UUID.String()).Set(ctx, CustomerToFirestore(customer))
	return err
}

func (fcr *FirestoreCustomerRepo) Delete(ctx context.Context, uuid uuid.UUID) error {
	_, err := fcr.client.Collection(customersCollectionName).Doc(uuid.String()).Delete(ctx)
	return err
}
