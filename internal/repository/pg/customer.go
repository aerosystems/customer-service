package pg

import (
	"context"
	"errors"
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type CustomerRepo struct {
	db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) *CustomerRepo {
	return &CustomerRepo{
		db: db,
	}
}

type Customer struct {
	Id        int       `gorm:"primaryKey;unique;autoIncrement"`
	Uuid      uuid.UUID `gorm:"unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (c *Customer) ToModel() *models.Customer {
	return &models.Customer{
		Uuid:      c.Uuid,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func CustomerToPg(customer *models.Customer) *Customer {
	return &Customer{
		Uuid:      customer.Uuid,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}

func (r *CustomerRepo) GetByUuid(ctx context.Context, uuid uuid.UUID) (*models.Customer, error) {
	var user Customer
	result := r.db.Find(&user, "uuid = ?", uuid.String())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return user.ToModel(), nil
}

func (r *CustomerRepo) Create(ctx context.Context, user *models.Customer) error {
	result := r.db.Create(CustomerToPg(user))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CustomerRepo) Update(ctx context.Context, user *models.Customer) error {
	result := r.db.Save(CustomerToPg(user))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CustomerRepo) Delete(ctx context.Context, uuid uuid.UUID) error {
	result := r.db.Delete(&Customer{}, "uuid = ?", uuid.String())
	if result.Error != nil {
		return result.Error
	}
	return nil
}
