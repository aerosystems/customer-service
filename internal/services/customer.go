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
	GetUserById(userId int) (*models.Customer, error)
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

func (us *CustomerServiceImpl) GetUserById(userId int) (*models.Customer, error) {
	user, err := us.customerRepo.GetById(userId)
	if err != nil {
		return nil, errors.New("could not get user id")
	}
	if user == nil {
		return nil, errors.New("user with this id does not exist")
	}
	return user, nil
}

func (us *CustomerServiceImpl) CreateUser() (user *models.Customer, err error) {
	defer func() {
		if r := recover(); r != nil {
			_ = us.customerRepo.Delete(user)
			_ = us.subsRPC.DeleteSubscription(user.Id)
			user = nil
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()
	user = NewCustomer()
	if err := us.customerRepo.Create(user); err != nil {
		return nil, errors.New("could not create new user")
	}
	if err := us.subsRPC.CreateFreeTrial(user.Id); err != nil {
		panic(errors.New("could not create free trial"))
	}
	if err := us.projectRPC.CreateDefaultProject(user.Id); err != nil {
		panic(errors.New("could not create default project"))
	}
	return
}
