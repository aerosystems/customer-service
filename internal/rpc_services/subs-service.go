package RPCServices

import (
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/google/uuid"
	"net/rpc"
)

type SubscriptionService interface {
	CreateFreeTrial(userId int) error
	DeleteSubscription(userId int) error
}

type SubscriptionRPC struct {
	rpcClient *rpc.Client
}

func NewSubsRPC(rpcClient *rpc.Client) *SubscriptionRPC {
	return &SubscriptionRPC{
		rpcClient: rpcClient,
	}
}

type SubscriptionRPCPayload struct {
	UserUuid uuid.UUID
	Kind     string
}

func (ss *SubscriptionRPC) CreateFreeTrial(customer *models.Customer) error {
	var resSub string
	err := ss.rpcClient.Call("SubsServer.CreateFreeTrial", SubscriptionRPCPayload{
		UserUuid: customer.Uuid,
		Kind:     "startup",
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
