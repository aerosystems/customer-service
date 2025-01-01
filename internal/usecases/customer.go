package usecases

import (
	"context"
	"github.com/aerosystems/customer-service/internal/domain"
	"github.com/sirupsen/logrus"
)

type CustomerUsecase struct {
	log                       *logrus.Logger
	customerRepo              CustomerRepository
	subscriptionEventsAdapter SubscriptionEventsAdapter
}

func NewCustomerUsecase(
	log *logrus.Logger,
	customerRepo CustomerRepository,
	subscriptionEventsAdapter SubscriptionEventsAdapter,
) *CustomerUsecase {
	return &CustomerUsecase{
		log:                       log,
		customerRepo:              customerRepo,
		subscriptionEventsAdapter: subscriptionEventsAdapter,
	}
}

func (cu CustomerUsecase) CreateCustomer(ctx context.Context, email, firebaseUID string) error {
	customer := domain.NewCustomer(email, firebaseUID)
	if err := cu.customerRepo.Create(ctx, customer); err != nil {
		return err
	}
	return nil
}
