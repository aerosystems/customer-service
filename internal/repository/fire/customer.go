package fire

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"github.com/aerosystems/customer-service/internal/models"
	"google.golang.org/api/iterator"
)

type CustomerRepo struct {
	client *firestore.Client
}

func NewCustomerRepo(client *firestore.Client) *CustomerRepo {
	return &CustomerRepo{
		client: client,
	}
}

func (r *CustomerRepo) GetAll(ctx context.Context) ([]models.Customer, error) {
	var customers []models.Customer

	iter := r.client.Collection("customers").Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		var customer models.Customer
		if err := doc.DataTo(&customer); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

func (r *CustomerRepo) GetById(ctx context.Context, id string) (*models.Customer, error) {
	doc, err := r.client.Collection("customers").Doc(id).Get(ctx)
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

func (r *CustomerRepo) Delete(ctx context.Context, id string) error {
	_, err := r.client.Collection("customers").Doc(id).Delete(ctx)
	return err
}
