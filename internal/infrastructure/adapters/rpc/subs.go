package RpcRepo

import (
	"github.com/aerosystems/customer-service/internal/models"
	RpcClient "github.com/aerosystems/customer-service/pkg/rpc_client"
	"github.com/google/uuid"
	"time"
)

type SubsRepo struct {
	rpcClient *RpcClient.ReconnectRpcClient
}

func NewSubsRepo(rpcClient *RpcClient.ReconnectRpcClient) *SubsRepo {
	return &SubsRepo{
		rpcClient: rpcClient,
	}
}

type SubsRPCPayload struct {
	UserUuid   uuid.UUID
	Kind       string
	AccessTime time.Time
}

func (ss *SubsRepo) CreateFreeTrial(customer *models.Customer) error {
	var resSub string
	err := ss.rpcClient.Call("Server.CreateFreeTrial", SubsRPCPayload{
		UserUuid: customer.Uuid,
		Kind:     models.TrialSubscription.String(),
	}, &resSub)
	if err != nil {
		return err
	}
	return nil
}

func (ss *SubsRepo) DeleteSubscription(customer *models.Customer) error {
	var resSub string
	err := ss.rpcClient.Call("Server.DeleteSubscription", customer.Uuid, &resSub)
	if err != nil {
		return err
	}
	return nil
}
