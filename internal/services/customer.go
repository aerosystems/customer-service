package services

import (
	"errors"
	"github.com/aerosystems/customer-service/internal/models"
)

type CustomerService interface {
	GetUserById(userId uint) (*models.Customer, error)
	GetUserByEmail(email string) (*models.Customer, error)
	CreateUser(email, passwordHash string) (*models.Customer, error)
	ResetPassword(userId uint, passwordHash string) error
	ActivateUser(userId uint) error
	MatchPassword(email, passwordHash string) (bool, error)
}

type CustomerServiceImpl struct {
	customerRepo models.CustomerRepository
}

func NewUserServiceImpl(userRepo models.CustomerRepository) *CustomerServiceImpl {
	return &CustomerServiceImpl{
		customerRepo: userRepo,
	}
}

func (us *CustomerServiceImpl) GetUserById(userId uint) (*models.Customer, error) {
	user, err := us.customerRepo.GetById(userId)
	if err != nil {
		return nil, errors.New("could not get user id")
	}
	if user == nil {
		return nil, errors.New("user with this id does not exist")
	}
	return user, nil
}

func (us *CustomerServiceImpl) CreateUser(email, passwordHash string) (*models.Customer, error) {
	user := models.NewUserEntity(email, passwordHash, "user")
	if err := us.customerRepo.Create(user); err != nil {
		return nil, errors.New("could not create new user")
	}
	return user, nil
}
