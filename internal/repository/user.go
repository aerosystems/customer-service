package repository

import (
	"errors"
	"github.com/aerosystems/user-service/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) FindAll() (*[]models.User, error) {
	var users []models.User
	r.db.Find(&users)
	return &users, nil
}

func (r *UserRepo) FindByID(ID int) (*models.User, error) {
	var user models.User
	result := r.db.Find(&user, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepo) FindByEmail(Email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", Email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepo) FindByGoogleID(GoogleID string) (*models.User, error) {
	var user models.User
	result := r.db.Where("google_id = ?", GoogleID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepo) Create(user *models.User) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepo) Update(user *models.User) error {
	result := r.db.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepo) Delete(user *models.User) error {
	result := r.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ResetPassword is the method we will use to change a user's password.
func (r *UserRepo) ResetPassword(user *models.User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	result := r.db.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (r *UserRepo) PasswordMatches(user *models.User, plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
