package adapters

import (
	"context"

	"github.com/google/uuid"

	"github.com/aerosystems/common-service/clients/grpcclient"
	"github.com/aerosystems/common-service/gen/protobuf/subscription"
	"github.com/aerosystems/customer-service/internal/usecases"
)

type SubscriptionAdapter struct {
	client subscription.SubscriptionServiceClient
}

func NewSubscriptionAdapter(address string) (*SubscriptionAdapter, error) {
	conn, err := grpcclient.NewGRPCConn(address)
	if err != nil {
		return nil, err
	}
	return &SubscriptionAdapter{
		client: subscription.NewSubscriptionServiceClient(conn),
	}, nil
}

func (sa SubscriptionAdapter) CreateFreeTrialSubscription(ctx context.Context, customerUUID uuid.UUID) (*usecases.SubscriptionDTO, error) {
	resp, err := sa.client.CreateFreeTrialSubscription(ctx, &subscription.CreateFreeTrialSubscriptionRequest{
		CustomerUuid: customerUUID.String(),
	})
	if err != nil {
		return nil, err
	}
	return &usecases.SubscriptionDTO{
		UUID:        resp.SubscriptionUuid,
		Type:        resp.SubscriptionType,
		AccessTime:  resp.AccessTime.AsTime(),
		AccessCount: resp.AccessCount,
	}, nil
}

func (sa SubscriptionAdapter) DeleteSubscription(ctx context.Context, subscriptionUUID uuid.UUID) error {
	_, err := sa.client.DeleteSubscription(ctx, &subscription.DeleteSubscriptionRequest{
		SubscriptionUuid: subscriptionUUID.String(),
	})
	return err
}
