package handlers

import (
	"context"
)

type CustomerUsecase interface {
	CreateCustomer(ctx context.Context, email, firebaseUID string) error
}
