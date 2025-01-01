package FirestoreRepo

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	CustomErrors "github.com/aerosystems/customer-service/internal/common/custom_errors"
	"github.com/aerosystems/customer-service/internal/domain"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const (
	defaultTimeout          = 2 * time.Second
	customersCollectionName = "customers"
)

type CustomerRepo struct {
	client *firestore.Client
}

func NewCustomerRepo(client *firestore.Client) *CustomerRepo {
	return &CustomerRepo{
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

func (r *CustomerRepo) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Customer, error) {
	c, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	docRef := r.client.Collection(customersCollectionName).Doc(uuid.String())
	doc, err := docRef.Get(c)
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

func (r *CustomerRepo) Create(ctx context.Context, customer *domain.Customer) error {
	currentCustomer, err := r.GetByUUID(ctx, customer.UUID)
	if err != nil && !errors.Is(err, CustomErrors.ErrCustomerNotFound) {
		return err
	}
	if currentCustomer != nil {
		return CustomErrors.ErrCustomerAlreadyExists
	}
	c, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err = r.client.Collection(customersCollectionName).Doc(customer.UUID.String()).Set(c, CustomerToFirestore(customer))
	return err
}

func (r *CustomerRepo) Update(ctx context.Context, customer *domain.Customer) error {
	c, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err := r.client.Collection(customersCollectionName).Doc(customer.UUID.String()).Set(c, CustomerToFirestore(customer))
	return err
}

func (r *CustomerRepo) Delete(ctx context.Context, uuid uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err := r.client.Collection(customersCollectionName).Doc(uuid.String()).Delete(c)
	return err
}
