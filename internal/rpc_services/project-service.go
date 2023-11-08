package RPCServices

import (
	"net/rpc"
)

type ProjectService interface {
	CreateDefaultProject(userId int) error
}

type ProjectRPC struct {
	rpcClient *rpc.Client
}

func NewProjectRPC(rpcClient *rpc.Client) *ProjectRPC {
	return &ProjectRPC{
		rpcClient: rpcClient,
	}
}

type ProjectRPCPayload struct {
	ID     int
	UserId int
	Name   string
	Token  string
}

func (ps *ProjectRPC) CreateDefaultProject(userId int) error {
	if err := ps.rpcClient.Call("ProjectServer.CreateProject", ProjectRPCPayload{
		UserId: userId,
		Name:   "default",
	}, nil); err != nil {
		return err
	}
	return nil
}
