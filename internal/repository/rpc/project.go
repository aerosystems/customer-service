package RpcRepo

import (
	"github.com/aerosystems/customer-service/internal/models"
	RpcClient "github.com/aerosystems/customer-service/pkg/rpc_client"
	"github.com/google/uuid"
)

type ProjectRepo struct {
	rpcClient *RpcClient.ReconnectRpcClient
}

func NewProjectRepo(rpcClient *RpcClient.ReconnectRpcClient) *ProjectRepo {
	return &ProjectRepo{
		rpcClient: rpcClient,
	}
}

type ProjectRPCPayload struct {
	Id       int
	UserUuid uuid.UUID
	Name     string
	Token    string
}

func (pr *ProjectRepo) CreateDefaultProject(customer *models.Customer) error {
	if err := pr.rpcClient.Call("Server.CreateDefaultProject", ProjectRPCPayload{
		UserUuid: customer.Uuid,
	}, nil); err != nil {
		return err
	}
	return nil
}
