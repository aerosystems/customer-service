package usecases

import (
	"errors"
	"fmt"
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

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

func NewCustomer() *models.Customer {
	return &models.Customer{
		Uuid: uuid.New(),
	}
}

func (cu *CustomerUsecase) GetUserByUuid(userUuid string) (*models.Customer, error) {
	uuid := uuid.MustParse(userUuid)
	user, err := cu.customerRepo.GetByUuid(uuid)
	if err != nil {
		return nil, errors.New("could not get user id")
	}
	if user == nil {
		return nil, errors.New("user with this id does not exist")
	}
	return user, nil
}

func (cu *CustomerUsecase) CreateUser() (customer *models.Customer, err error) {
	defer func() {
		if r := recover(); r != nil {
			_ = cu.customerRepo.Delete(customer)
			_ = cu.subsRepo.DeleteSubscription(customer)
			customer = nil
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()
	customer = NewCustomer()
	if err := cu.customerRepo.Create(customer); err != nil {
		return nil, errors.New("could not create new customer")
	}
	if err := cu.subsRepo.CreateFreeTrial(customer); err != nil {
		panic(errors.New("could not create free trial"))
	}
	if err := cu.projectRepo.CreateDefaultProject(customer); err != nil {
		log.Errorf("could not create default project: %v", err)
		panic(errors.New("could not create default project"))
	}
	return
}
