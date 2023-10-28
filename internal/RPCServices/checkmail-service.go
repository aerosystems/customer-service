package RPCServices

import (
	"errors"
	"net/rpc"
)

type CheckmailService interface {
	IsTrustEmail(email, clientIp string) (bool, error)
}

type CheckmailRPC struct {
	rpcClient *rpc.Client
}

type InspectRPCPayload struct {
	Domain   string
	ClientIp string
}

func NewCheckmailRPC(rpcClient *rpc.Client) *CheckmailRPC {
	return &CheckmailRPC{
		rpcClient: rpcClient,
	}
}

func (cs *CheckmailRPC) IsTrustEmail(email, clientIp string) (bool, error) {
	var result string
	if err := cs.rpcClient.Call(
		"CheckmailServer.Inspect",
		InspectRPCPayload{
			Domain:   email,
			ClientIp: clientIp,
		},
		&result); err != nil {
		return false, errors.New("email address does not valid")
	}

	if result == "blacklist" {
		return false, errors.New("email address contains in blacklist")
	}

	return true, nil
}
