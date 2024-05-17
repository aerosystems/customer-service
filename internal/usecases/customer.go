package usecases

import (
	"context"
	"fmt"
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/google/uuid"
	"time"
)

const defaultTimeout = 2 * time.Second

type CustomerUsecase struct {
	customerRepo CustomerRepository
	projectRepo  ProjectRepository
	subsRepo     SubsRepository
}

func NewCustomerUsecase(
	customerRepo CustomerRepository,
	projectRepo ProjectRepository,
	subsRepo SubsRepository,
) *CustomerUsecase {
	return &CustomerUsecase{
		customerRepo: customerRepo,
		projectRepo:  projectRepo,
		subsRepo:     subsRepo,
	}
}

func (cu *CustomerUsecase) CreateCustomer(uuidStr string) (customer *models.Customer, err error) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		ctx := context.Background()
	//		_ = cu.customerRepo.Delete(ctx, customer.Uuid)
	//		_ = cu.subsRepo.DeleteSubscription(customer)
	//		customer = nil
	//		err = fmt.Errorf("panic occurred: %v", r)
	//	}
	//}()
	customerUuid, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, fmt.Errorf("could not parse uuid: %v", err)
	}
	customer = &models.Customer{
		Uuid:      customerUuid,
		CreatedAt: time.Now(),
	}
	ctx := context.Background()
	if err := cu.customerRepo.Create(ctx, customer); err != nil {
		return nil, err
	}
	//if err := cu.subsRepo.CreateFreeTrial(customer); err != nil {
	//	log.Errorf("could not create free trial: %v", err)
	//	panic(errors.New("could not create free trial"))
	//}
	//if err := cu.projectRepo.CreateDefaultProject(customer); err != nil {
	//	log.Errorf("could not create default project: %v", err)
	//	panic(errors.New("could not create default project"))
	//}
	return
}
