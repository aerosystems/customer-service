package repository

import (
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepo struct {
	db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) *CustomerRepo {
	return &CustomerRepo{
		db: db,
	}
}

func (r *CustomerRepo) GetAll() (*[]models.Customer, error) {
	var users []models.Customer
	r.db.Find(&users)
	return &users, nil
}

func (r *CustomerRepo) GetById(Id int) (*models.Customer, error) {
	var user models.Customer
	result := r.db.Find(&user, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *CustomerRepo) GetByUuid(uuid uuid.UUID) (*models.Customer, error) {
	var user models.Customer
	result := r.db.Find(&user, "uuid = ?", uuid.String())
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *CustomerRepo) Create(user *models.Customer) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CustomerRepo) Update(user *models.Customer) error {
	result := r.db.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CustomerRepo) Delete(user *models.Customer) error {
	result := r.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
