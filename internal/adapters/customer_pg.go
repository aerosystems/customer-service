package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aerosystems/customer-service/internal/entities"
)

type CustomerPostgresRepo struct {
	db *gorm.DB
}

func NewCustomerPostgresRepo(db *gorm.DB) *CustomerPostgresRepo {
	return &CustomerPostgresRepo{
		db: db,
	}
}

type Customer struct {
	UUID        string
	Email       string
	FirebaseUID string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	DeleteAt    *time.Time
}

func (c *Customer) ToModel() *entities.Customer {
	return &entities.Customer{
		UUID:        uuid.MustParse(c.UUID),
		Email:       c.Email,
		FirebaseUID: c.FirebaseUID,
		CreatedAt:   c.CreatedAt,
	}
}

func CustomerToPostgres(customer *entities.Customer) *Customer {
	return &Customer{
		UUID:        customer.UUID.String(),
		Email:       customer.Email,
		FirebaseUID: customer.FirebaseUID,
		CreatedAt:   customer.CreatedAt,
		DeleteAt:    nil,
	}
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.UUID == "" {
		c.UUID = uuid.New().String()
	}
	return nil
}

func (c *Customer) BeforeDelete(tx *gorm.DB) error {
	if c.DeleteAt == nil {
		now := time.Now()
		c.DeleteAt = &now
	}
	return nil
}

func (cr *CustomerPostgresRepo) GetByCustomerUUID(ctx context.Context, customerUUID uuid.UUID) (*entities.Customer, error) {
	var customerPG Customer
	result := cr.db.WithContext(ctx).First(&customerPG, "uuid = ?", customerUUID.String())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, entities.ErrCustomerNotFound
		}
		return nil, result.Error
	}
	return customerPG.ToModel(), nil
}

func (cr *CustomerPostgresRepo) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*entities.Customer, error) {
	var customerPG Customer
	result := cr.db.WithContext(ctx).First(&customerPG, "firebase_uid = ?", firebaseUID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, entities.ErrCustomerNotFound
		}
		return nil, result.Error
	}
	return customerPG.ToModel(), nil
}

func (cr *CustomerPostgresRepo) Create(ctx context.Context, customer *entities.Customer) error {
	customerPG := CustomerToPostgres(customer)
	result := cr.db.WithContext(ctx).Create(customerPG)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return entities.ErrCustomerAlreadyExists
		}
		return result.Error
	}
	return nil
}

func (cr *CustomerPostgresRepo) Upsert(ctx context.Context, customer *entities.Customer) error {
	customerPG := CustomerToPostgres(customer)
	result := cr.db.WithContext(ctx).Save(customerPG)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return entities.ErrCustomerAlreadyExists
		}
		return result.Error
	}
	return nil
}

func (cr *CustomerPostgresRepo) Delete(ctx context.Context, customerUUID uuid.UUID) error {
	customerPG := &Customer{UUID: customerUUID.String()}
	result := cr.db.WithContext(ctx).Delete(customerPG)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
