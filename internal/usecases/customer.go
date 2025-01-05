package usecases

import (
	"context"
	"errors"
	CustomErrors "github.com/aerosystems/customer-service/internal/common/custom_errors"
	"github.com/aerosystems/customer-service/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CustomerUsecase struct {
	log                 *logrus.Logger
	customerRepo        CustomerRepository
	subscriptionAdapter SubscriptionAdapter
	projectAdapter      ProjectAdapter
}

func NewCustomerUsecase(
	log *logrus.Logger,
	customerRepo CustomerRepository,
	subscriptionAdapter SubscriptionAdapter,
	projectAdapter ProjectAdapter,
) *CustomerUsecase {
	return &CustomerUsecase{
		log:                 log,
		customerRepo:        customerRepo,
		subscriptionAdapter: subscriptionAdapter,
		projectAdapter:      projectAdapter,
	}
}

func (cu CustomerUsecase) CreateCustomer(ctx context.Context, email, firebaseUID string) error {
	customer, err := cu.customerRepo.GetByFirebaseUID(ctx, firebaseUID)
	if err != nil && !errors.Is(err, CustomErrors.ErrCustomerNotFound) {
		return err
	}
	if customer != nil {
		return CustomErrors.ErrCustomerAlreadyExists
	}
	customer = domain.NewCustomer(email, firebaseUID)
	var subscriptionUUID, projectUUID uuid.UUID
	defer func() {
		if err != nil {
			cu.compensationCreateCustomerError(ctx, err, subscriptionUUID, projectUUID)
		}
	}()
	//if subscriptionUUID, err = cu.subscriptionAdapter.CreateFreeTrialSubscription(ctx, customer.UUID); err != nil {
	//	return err
	//}
	if projectUUID, err = cu.projectAdapter.CreateDefaultProject(ctx, customer.UUID); err != nil {
		return err
	}
	if err = cu.customerRepo.Create(ctx, customer); err != nil {
		return err
	}
	return nil
}

func (cu CustomerUsecase) compensationCreateCustomerError(ctx context.Context, err error, subscriptionUUID, projectUUID uuid.UUID) {
	cu.log.WithError(err).Error("Failed creating customer error: %v", err)
	if subscriptionUUID != uuid.Nil {
		if compErr := cu.subscriptionAdapter.DeleteSubscription(ctx, subscriptionUUID); err != nil {
			cu.log.WithError(compErr).Error("Failed compensation creating subscription error: %v", err)
		}
	}
	if projectUUID != uuid.Nil {
		if compErr := cu.projectAdapter.DeleteProject(ctx, projectUUID); err != nil {
			cu.log.WithError(compErr).Error("Failed compensation creating project error: %v", err)
		}
	}
}
