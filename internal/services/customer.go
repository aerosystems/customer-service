package services

import (
	"errors"
	"fmt"
	"github.com/aerosystems/customer-service/internal/models"
	RPCServices "github.com/aerosystems/customer-service/internal/rpc_services"
	"github.com/google/uuid"
)

type CustomerService interface {
	CreateUser() (*models.Customer, error)
	GetUserByUuid(userUuid string) (*models.Customer, error)
}

type CustomerServiceImpl struct {
	customerRepo models.CustomerRepository
	projectRPC   *RPCServices.ProjectRPC
	subsRPC      *RPCServices.SubscriptionRPC
}

func NewCustomerServiceImpl(customerRepository models.CustomerRepository, projectRPC *RPCServices.ProjectRPC, subsRPC *RPCServices.SubscriptionRPC) *CustomerServiceImpl {
	return &CustomerServiceImpl{
		customerRepo: customerRepository,
		projectRPC:   projectRPC,
		subsRPC:      subsRPC,
	}
}

func NewCustomer() *models.Customer {
	return &models.Customer{
		Uuid: uuid.New(),
	}
}

func (us *CustomerServiceImpl) GetUserByUuid(userUuid string) (*models.Customer, error) {
	uuid := uuid.MustParse(userUuid)
	user, err := us.customerRepo.GetByUuid(uuid)
	if err != nil {
		return nil, errors.New("could not get user id")
	}
	if user == nil {
		return nil, errors.New("user with this id does not exist")
	}
	return user, nil
}

func (us *CustomerServiceImpl) CreateUser() (customer *models.Customer, err error) {
	defer func() {
		if r := recover(); r != nil {
			_ = us.customerRepo.Delete(customer)
			_ = us.subsRPC.DeleteSubscription(customer)
			customer = nil
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()
	customer = NewCustomer()
	if err := us.customerRepo.Create(customer); err != nil {
		return nil, errors.New("could not create new customer")
	}
	if err := us.subsRPC.CreateFreeTrial(customer); err != nil {
		panic(errors.New("could not create free trial"))
	}
	if err := us.projectRPC.CreateDefaultProject(customer); err != nil {
		panic(errors.New("could not create default project"))
	}
	return
}
