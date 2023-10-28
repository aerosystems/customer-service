package RPCServices

import (
	"net/rpc"
)

type SubscriptionService interface {
	CreateFreeTrial(userId uint) error
}

type SubscriptionRPC struct {
	rpcClient *rpc.Client
}

func NewSubscriptionRPC(rpcClient *rpc.Client) *SubscriptionRPC {
	return &SubscriptionRPC{
		rpcClient: rpcClient,
	}
}

type SubscriptionRPCPayload struct {
	UserId uint
	Kind   string
}

func (ss *SubscriptionRPC) CreateFreeTrial(userId uint) error {
	var resSub string
	err := ss.rpcClient.Call("SubsServer.CreateFreeTrial", SubscriptionRPCPayload{
		UserId: userId,
		Kind:   "startup",
	}, &resSub)
	if err != nil {
		return err
	}
	return nil
}
