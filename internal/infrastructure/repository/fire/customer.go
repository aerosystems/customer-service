package fire

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	CustomErrors "github.com/aerosystems/customer-service/internal/common/custom_errors"
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const (
	defaultTimeout = 2 * time.Second
)

type CustomerRepo struct {
	client *firestore.Client
}

func NewCustomerRepo(client *firestore.Client) *CustomerRepo {
	return &CustomerRepo{
		client: client,
	}
}

type CustomerFirestore struct {
	Uuid      string    `firestore:"uuid"`
	CreatedAt time.Time `firestore:"created_at"`
}

func (c *CustomerFirestore) ToModel() *models.Customer {
	return &models.Customer{
		Uuid:      uuid.MustParse(c.Uuid),
		CreatedAt: c.CreatedAt,
	}
}

func CustomerToFirestore(customer *models.Customer) *CustomerFirestore {
	return &CustomerFirestore{
		Uuid:      customer.Uuid.String(),
		CreatedAt: customer.CreatedAt,
	}
}

func (r *CustomerRepo) GetByUuid(ctx context.Context, uuid uuid.UUID) (*models.Customer, error) {
	c, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	docRef := r.client.Collection("customers").Doc(uuid.String())
	doc, err := docRef.Get(c)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, CustomErrors.ErrCustomerNotFound
		}
		return nil, err
	}

	var customer CustomerFirestore
	if err := doc.DataTo(&customer); err != nil {
		return nil, err
	}

	return customer.ToModel(), nil
}

func (r *CustomerRepo) Create(ctx context.Context, customer *models.Customer) error {
	currentCustomer, err := r.GetByUuid(ctx, customer.Uuid)
	if err != nil && !errors.Is(err, CustomErrors.ErrCustomerNotFound) {
		return err
	}
	if currentCustomer != nil {
		return CustomErrors.ErrCustomerAlreadyExists
	}
	c, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	_, err = r.client.Collection("customers").Doc(customer.Uuid.String()).Set(c, CustomerToFirestore(customer))
	return err
}

func (r *CustomerRepo) Update(ctx context.Context, customer *models.Customer) error {
	c, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	_, err := r.client.Collection("customers").Doc(customer.Uuid.String()).Set(c, CustomerToFirestore(customer))
	return err
}

func (r *CustomerRepo) Delete(ctx context.Context, uuid uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	_, err := r.client.Collection("customers").Doc(uuid.String()).Delete(c)
	return err
}
