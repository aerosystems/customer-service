package services

import (
	"errors"
	"github.com/aerosystems/user-service/internal/models"
	"gorm.io/gorm"
)

type UserService interface {
	GetUser(email string) (*models.User, error)
	CreateUser(email, password string) (uint, error)
	ResetPassword(email, password string) error
}

type UserServiceImpl struct {
	userRepo models.UserRepository
}

func NewUserServiceImpl(userRepo models.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (us *UserServiceImpl) GetUser(email string) (*models.User, error) {
	user, err := us.userRepo.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("could not get user id")
	}
	if user == nil {
		return nil, errors.New("user with this email does not exist")
	}
	return user, nil
}

func (us *UserServiceImpl) CreateUser(email, password string) (uint, error) {
	newUser := models.NewUser(email, password, "user")
	if err := us.userRepo.Create(newUser); err != nil {
		return 0, errors.New("could not create new user")
	}
	return newUser.Id, nil
}

func (us *UserServiceImpl) ResetPassword(email, password string) error {
	user, err := us.userRepo.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("could not get user id")
	}
	if user == nil {
		return errors.New("user with this email does not exist")
	}
	if err := us.userRepo.ResetPassword(user, password); err != nil {
		return errors.New("could not update password")
	}
	return nil
}
