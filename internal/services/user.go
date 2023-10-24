package services

import (
	"errors"
	"github.com/aerosystems/user-service/internal/models"
)

type UserService interface {
	GetUserById(userId uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(email, passwordHash string) (*models.User, error)
	ResetPassword(userId uint, passwordHash string) error
	ActivateUser(userId uint) error
	MatchPassword(email, passwordHash string) (bool, error)
}

type UserServiceImpl struct {
	userRepo models.UserRepository
}

func NewUserServiceImpl(userRepo models.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (us *UserServiceImpl) GetUserById(userId uint) (*models.User, error) {
	user, err := us.userRepo.FindById(userId)
	if err != nil {
		return nil, errors.New("could not get user id")
	}
	if user == nil {
		return nil, errors.New("user with this id does not exist")
	}
	return user, nil
}

func (us *UserServiceImpl) GetUserByEmail(email string) (*models.User, error) {
	user, err := us.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("could not get user id")
	}
	if user == nil {
		return nil, errors.New("user with this email does not exist")
	}
	return user, nil
}

func (us *UserServiceImpl) CreateUser(email, passwordHash string) (*models.User, error) {
	user := models.NewUserEntity(email, passwordHash, "user")
	if err := us.userRepo.Create(user); err != nil {
		return nil, errors.New("could not create new user")
	}
	return user, nil
}

func (us *UserServiceImpl) ResetPassword(userId uint, passwordHash string) error {
	user, err := us.userRepo.FindById(userId)
	if err != nil {
		return errors.New("could not get user id")
	}
	if user == nil {
		return errors.New("user with this email does not exist")
	}
	user.PasswordHash = passwordHash
	if err := us.userRepo.Update(user); err != nil {
		return errors.New("could not update password")
	}
	return nil
}

func (us *UserServiceImpl) ActivateUser(userId uint) error {
	user, err := us.userRepo.FindById(userId)
	if err != nil {
		return errors.New("could not get user id")
	}
	if user == nil {
		return errors.New("user with this email does not exist")
	}
	user.IsActive = true
	if err := us.userRepo.Update(user); err != nil {
		return errors.New("could not activate user")
	}
	return nil
}

func (us *UserServiceImpl) MatchPassword(email, passwordHash string) (bool, error) {
	user, err := us.userRepo.FindByEmail(email)
	if err != nil {
		return false, errors.New("could not get user id")
	}
	if user == nil {
		return false, errors.New("user with this email does not exist")
	}
	if user.PasswordHash != passwordHash {
		return false, errors.New("password is incorrect")
	}
	return true, nil
}
