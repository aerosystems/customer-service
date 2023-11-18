package RPCServices

import (
	"github.com/aerosystems/customer-service/internal/models"
	RPCClient "github.com/aerosystems/customer-service/pkg/rpc_client"
	"github.com/google/uuid"
)

type ProjectService interface {
	CreateDefaultProject(userId int) error
}

type ProjectRPC struct {
	rpcClient *RPCClient.ReconnectRPCClient
}

func NewProjectRPC(rpcClient *RPCClient.ReconnectRPCClient) *ProjectRPC {
	return &ProjectRPC{
		rpcClient: rpcClient,
	}
}

type ProjectRPCPayload struct {
	Id       int
	UserUuid uuid.UUID
	Name     string
	Token    string
}

func (ps *ProjectRPC) CreateDefaultProject(customer *models.Customer) error {
	if err := ps.rpcClient.Call("ProjectServer.CreateDefaultProject", ProjectRPCPayload{
		UserUuid: customer.Uuid,
	}, nil); err != nil {
		return err
	}
	return nil
}
