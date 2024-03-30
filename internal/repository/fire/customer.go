package fire

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/google/uuid"
)

type CustomerRepo struct {
	client *firestore.Client
	ctx    context.Context
}

func NewCustomerRepo(client *firestore.Client, ctx context.Context) *CustomerRepo {
	return &CustomerRepo{
		client: client,
		ctx:    ctx,
	}
}

func (r *CustomerRepo) GetByUuid(uuid uuid.UUID) (*models.Customer, error) {
	docRef := r.client.Collection("customers").Doc(uuid.String())
	doc, err := docRef.Get(r.ctx)
	if err != nil {
		return nil, err
	}

	var customer models.Customer
	if err := doc.DataTo(&customer); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerRepo) Create(customer *models.Customer) error {
	_, err := r.client.Collection("customers").Doc(customer.Uuid.String()).Set(r.ctx, customer)
	return err
}

func (r *CustomerRepo) Update(customer *models.Customer) error {
	_, err := r.client.Collection("customers").Doc(customer.Uuid.String()).Set(r.ctx, customer)
	return err
}

func (r *CustomerRepo) Delete(customer *models.Customer) error {
	_, err := r.client.Collection("customers").Doc(customer.Uuid.String()).Delete(r.ctx)
	return err
}
