package RpcRepo

import (
	"github.com/aerosystems/customer-service/internal/models"
	RpcClient "github.com/aerosystems/customer-service/pkg/rpc_client"
)

type SubsRepo struct {
	rpcClient *RpcClient.ReconnectRpcClient
}

func NewSubsRepo(rpcClient *RpcClient.ReconnectRpcClient) *SubsRepo {
	return &SubsRepo{
		rpcClient: rpcClient,
	}
}

func (ss *SubsRepo) CreateFreeTrial(customer *models.Customer) error {
	var resSub string
	err := ss.rpcClient.Call("Server.CreateFreeTrial", models.SubsRPCPayload{
		UserUuid: customer.Uuid,
		Kind:     models.TrialSubscription,
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
