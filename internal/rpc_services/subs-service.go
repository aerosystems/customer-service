package RPCServices

import (
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
	UserId int
	Kind   string
}

func (ss *SubscriptionRPC) CreateFreeTrial(userId int) error {
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

func (ss *SubscriptionRPC) DeleteSubscription(userId int) error {
	var resSub string
	err := ss.rpcClient.Call("SubsServer.DeleteSubscription", userId, &resSub)
	if err != nil {
		return err
	}
	return nil
}
