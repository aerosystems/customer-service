package RPCServices

import (
	"net/rpc"
)

type ProjectService interface {
	CreateDefaultProject(userId uint) error
}

type ProjectRPC struct {
	rpcClient *rpc.Client
}

func NewProjectRPC(rpcClient *rpc.Client) *ProjectRPC {
	return &ProjectRPC{
		rpcClient: rpcClient,
	}
}

type CreateProjectRPCPayload struct {
	UserId uint
	Name   string
}

func (ps *ProjectRPC) CreateDefaultProject(userId uint) error {
	var result string
	if err := ps.rpcClient.Call("ProjectServer.CreateProject", CreateProjectRPCPayload{
		UserId: userId,
		Name:   "default",
	}, &result); err != nil {
		return err
	}
	return nil
}
