package usecases

import (
	"context"
	"fmt"
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
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

func (cu CustomerUsecase) CreateCustomer(uuidStr string) (customer *models.Customer, err error) {
	customerUuid, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("could not parse uuid: %w", err)
	}
	customer = &models.Customer{
		Uuid:      customerUuid,
		CreatedAt: time.Now(),
	}
	ctx := context.Background()
	if err := cu.customerRepo.Create(ctx, customer); err != nil {
		return nil, err
	}
	if err := cu.subscriptionEventsAdapter.PublishCreateSubscriptionEvent(
		customerUuid,
		models.TrialSubscription,
		models.OneWeekSubscriptionDuration,
	); err != nil {
		cu.log.Errorf("could not publish create subscription event: %v", err)
	}
	return customer, nil
}
