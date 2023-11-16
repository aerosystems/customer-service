package RPCServices

import (
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/google/uuid"
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
	Id       int
	UserUuid uuid.UUID
	Name     string
	Token    string
}

func (ps *ProjectRPC) CreateDefaultProject(customer *models.Customer) error {
	if err := ps.rpcClient.Call("ProjectServer.CreateProject", ProjectRPCPayload{
		UserUuid: customer.Uuid,
		Name:     "default",
	}, nil); err != nil {
		return err
	}
	return nil
}
