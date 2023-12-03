package RPCServices

import (
	"github.com/aerosystems/customer-service/internal/models"
	RPCClient "github.com/aerosystems/customer-service/pkg/rpc_client"
	"github.com/google/uuid"
	"time"
)

type SubscriptionService interface {
	CreateFreeTrial(userId int) error
	DeleteSubscription(userId int) error
}

type SubscriptionRPC struct {
	rpcClient *RPCClient.ReconnectRPCClient
}

func NewSubsRPC(rpcClient *RPCClient.ReconnectRPCClient) *SubscriptionRPC {
	return &SubscriptionRPC{
		rpcClient: rpcClient,
	}
}

type SubsRPCPayload struct {
	UserUuid   uuid.UUID
	Kind       models.KindSubscription
	AccessTime time.Time
}

func (ss *SubscriptionRPC) CreateFreeTrial(customer *models.Customer) error {
	var resSub string
	err := ss.rpcClient.Call("SubsServer.CreateFreeTrial", SubsRPCPayload{
		UserUuid: customer.Uuid,
		Kind:     models.TrialSubscription,
	}, &resSub)
	if err != nil {
		return err
	}
	return nil
}

func (ss *SubscriptionRPC) DeleteSubscription(customer *models.Customer) error {
	var resSub string
	err := ss.rpcClient.Call("SubsServer.DeleteSubscription", customer.Uuid, &resSub)
	if err != nil {
		return err
	}
	return nil
}
