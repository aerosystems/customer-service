package adapters

import (
	"context"
	"github.com/google/uuid"
)

type SubscriptionAdapter struct {
}

func NewSubscriptionAdapter(address string) (*SubscriptionAdapter, error) {
	return &SubscriptionAdapter{}, nil
}

func (sa SubscriptionAdapter) CreateFreeTrialSubscription(ctx context.Context, customerUUID uuid.UUID) (subscriptionUUID uuid.UUID, err error) {
	return uuid.Nil, nil
}
func (sa SubscriptionAdapter) DeleteSubscription(ctx context.Context, subscriptionUUID uuid.UUID) error {
	return nil
}
