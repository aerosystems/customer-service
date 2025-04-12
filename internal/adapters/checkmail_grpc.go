package adapters

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/aerosystems/common-service/clients/grpcclient"
	"github.com/aerosystems/common-service/gen/protobuf/checkmail"
	"github.com/aerosystems/customer-service/internal/usecases"
)

type CheckmailAdapter struct {
	client checkmail.CheckmailServiceClient
}

func NewCheckmailAdapter(address string) (*CheckmailAdapter, error) {
	conn, err := grpcclient.NewGRPCConn(address)
	if err != nil {
		return nil, err
	}
	return &CheckmailAdapter{
		client: checkmail.NewCheckmailServiceClient(conn),
	}, nil
}

func (ca CheckmailAdapter) CreateAccess(ctx context.Context, projectToken string, subscriptionDTO *usecases.SubscriptionDTO) error {
	_, err := ca.client.CreateAccess(ctx, &checkmail.CreateAccessRequest{
		ProjectToken:     projectToken,
		SubscriptionType: subscriptionDTO.Type,
		AccessTime:       timestamppb.New(subscriptionDTO.AccessTime),
		AccessCount:      subscriptionDTO.AccessCount,
	})
	return err
}
