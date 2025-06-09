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
	UUID        string     `gorm:"column:uuid;primaryKey;type:varchar(36)"`
	Email       string     `gorm:"column:email;uniqueIndex;not null"`
	FirebaseUID string     `gorm:"column:firebase_uid;uniqueIndex;not null"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index"` // For soft deletes
}

func (*Customer) TableName() string {
	return "customers"
}

func (c *Customer) ToModel() *entities.Customer {
	return &entities.Customer{
		UUID:        uuid.MustParse(c.UUID),
		Email:       c.Email,
		FirebaseUID: c.FirebaseUID,
		CreatedAt:   c.CreatedAt,
	}
}

func ModelToCustomer(customer *entities.Customer) *Customer {
	return &Customer{
		UUID:        customer.UUID.String(),
		Email:       customer.Email,
		FirebaseUID: customer.FirebaseUID,
		CreatedAt:   customer.CreatedAt,
		UpdatedAt:   time.Now(), // Set current time for updates
		DeletedAt:   nil,
	}
}

// BeforeCreate hook - only set UUID if it's empty
func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.UUID == "" {
		c.UUID = uuid.New().String()
	}
	return nil
}

func (r *CustomerPostgresRepo) GetByCustomerUUID(ctx context.Context, customerUUID uuid.UUID) (*entities.Customer, error) {
	var customer Customer
	err := r.db.WithContext(ctx).
		Where("uuid = ? AND deleted_at IS NULL", customerUUID.String()).
		First(&customer).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrCustomerNotFound
		}
		return nil, err
	}

	return customer.ToModel(), nil
}

func (r *CustomerPostgresRepo) GetByFirebaseUID(ctx context.Context, firebaseUID string) (*entities.Customer, error) {
	if firebaseUID == "" {
		return nil, errors.New("firebase UID cannot be empty")
	}

	var customer Customer
	err := r.db.WithContext(ctx).
		Where("firebase_uid = ? AND deleted_at IS NULL", firebaseUID).
		First(&customer).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrCustomerNotFound
		}
		return nil, err
	}

	return customer.ToModel(), nil
}

func (r *CustomerPostgresRepo) GetByEmail(ctx context.Context, email string) (*entities.Customer, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	var customer Customer
	err := r.db.WithContext(ctx).
		Where("email = ? AND deleted_at IS NULL", email).
		First(&customer).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrCustomerNotFound
		}
		return nil, err
	}

	return customer.ToModel(), nil
}

func (r *CustomerPostgresRepo) Create(ctx context.Context, customer *entities.Customer) error {
	if customer == nil {
		return errors.New("customer cannot be nil")
	}

	pgCustomer := ModelToCustomer(customer)
	err := r.db.WithContext(ctx).Create(pgCustomer).Error

	if err != nil {
		// Check for unique constraint violations
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			isDuplicateKeyError(err) {
			return entities.ErrCustomerAlreadyExists
		}
		return err
	}

	return nil
}

func (r *CustomerPostgresRepo) Update(ctx context.Context, customer *entities.Customer) error {
	if customer == nil {
		return errors.New("customer cannot be nil")
	}

	pgCustomer := ModelToCustomer(customer)
	result := r.db.WithContext(ctx).
		Where("uuid = ? AND deleted_at IS NULL", customer.UUID.String()).
		Updates(pgCustomer)

	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return entities.ErrCustomerAlreadyExists
		}
		return result.Error
	}

	if result.RowsAffected == 0 {
		return entities.ErrCustomerNotFound
	}

	return nil
}

func (r *CustomerPostgresRepo) Upsert(ctx context.Context, customer *entities.Customer) error {
	if customer == nil {
		return errors.New("customer cannot be nil")
	}

	// Try to find existing customer first
	existing, err := r.GetByCustomerUUID(ctx, customer.UUID)
	if err != nil && !errors.Is(err, entities.ErrCustomerNotFound) {
		return err
	}

	if existing != nil {
		// Update existing customer
		return r.Update(ctx, customer)
	}

	// Create new customer
	return r.Create(ctx, customer)
}

func (r *CustomerPostgresRepo) Delete(ctx context.Context, customerUUID uuid.UUID) error {
	// Soft delete - set deleted_at timestamp
	result := r.db.WithContext(ctx).
		Model(&Customer{}).
		Where("uuid = ? AND deleted_at IS NULL", customerUUID.String()).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return entities.ErrCustomerNotFound
	}

	return nil
}

func (r *CustomerPostgresRepo) HardDelete(ctx context.Context, customerUUID uuid.UUID) error {
	// Permanently delete from database
	result := r.db.WithContext(ctx).
		Where("uuid = ?", customerUUID.String()).
		Delete(&Customer{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return entities.ErrCustomerNotFound
	}

	return nil
}

func (r *CustomerPostgresRepo) Restore(ctx context.Context, customerUUID uuid.UUID) error {
	// Restore soft-deleted customer
	result := r.db.WithContext(ctx).
		Model(&Customer{}).
		Where("uuid = ? AND deleted_at IS NOT NULL", customerUUID.String()).
		Update("deleted_at", nil)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return entities.ErrCustomerNotFound
	}

	return nil
}

func (r *CustomerPostgresRepo) List(ctx context.Context, limit, offset int) ([]entities.Customer, error) {
	var customers []Customer
	query := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&customers).Error
	if err != nil {
		return nil, err
	}

	result := make([]entities.Customer, 0, len(customers))
	for _, customer := range customers {
		result = append(result, *customer.ToModel())
	}

	return result, nil
}

func (r *CustomerPostgresRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Customer{}).
		Where("deleted_at IS NULL").
		Count(&count).Error

	return count, err
}

// Helper function to check for duplicate key errors across different PostgreSQL drivers
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	// Check for common PostgreSQL duplicate key error patterns
	return contains(errStr, "duplicate key") ||
		contains(errStr, "UNIQUE constraint") ||
		contains(errStr, "violates unique constraint") ||
		contains(errStr, "duplicate entry")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(hasPrefix(s, substr) || hasSuffix(s, substr) || containsSubstring(s, substr))))
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func containsSubstring(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
