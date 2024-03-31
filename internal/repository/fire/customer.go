package fire

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/google/uuid"
)

type CustomerRepo struct {
	client *firestore.Client
}

func NewCustomerRepo(client *firestore.Client) *CustomerRepo {
	return &CustomerRepo{
		client: client,
	}
}

func (r *CustomerRepo) GetByUuid(ctx context.Context, uuid uuid.UUID) (*models.Customer, error) {
	docRef := r.client.Collection("customers").Doc(uuid.String())
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	var customer models.Customer
	if err := doc.DataTo(&customer); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerRepo) Create(ctx context.Context, customer *models.Customer) error {
	_, err := r.client.Collection("customers").Doc(customer.Uuid.String()).Set(ctx, customer)
	return err
}

func (r *CustomerRepo) Update(ctx context.Context, customer *models.Customer) error {
	_, err := r.client.Collection("customers").Doc(customer.Uuid.String()).Set(ctx, customer)
	return err
}

func (r *CustomerRepo) Delete(ctx context.Context, customer *models.Customer) error {
	_, err := r.client.Collection("customers").Doc(customer.Uuid.String()).Delete(ctx)
	return err
}
