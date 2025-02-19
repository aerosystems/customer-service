package adapters

import (
	"context"
	"crypto/tls"
	"github.com/aerosystems/common-service/gen/protobuf/checkmail"
	"github.com/aerosystems/customer-service/internal/usecases"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CheckmailAdapter struct {
	client checkmail.CheckmailServiceClient
}

func NewCheckmailAdapter(address string) (*CheckmailAdapter, error) {
	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30,
			Timeout: 30,
		}),
	}
	if address[len(address)-4:] == ":443" {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.NewClient(address, opts...)
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
