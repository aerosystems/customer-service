package entities

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	UUID        uuid.UUID
	Email       string
	FirebaseUID string
	CreatedAt   time.Time
}

func NewCustomer(email, firebaseUID string) *Customer {
	return &Customer{
		UUID:        uuid.New(),
		Email:       email,
		FirebaseUID: firebaseUID,
		CreatedAt:   time.Now(),
	}
}
