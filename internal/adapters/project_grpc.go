package adapters

import (
	"context"
	"github.com/aerosystems/customer-service/internal/common/protobuf/project"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type ProjectAdapter struct {
	client project.ProjectServiceClient
}

func NewProjectAdapter(address string, opts ...grpc.DialOption) (*ProjectAdapter, error) {
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, err
	}
	return &ProjectAdapter{
		client: project.NewProjectServiceClient(conn),
	}, nil
}

func (pa ProjectAdapter) CreateDefaultProject(ctx context.Context, customerUUID uuid.UUID) (projectUUID uuid.UUID, err error) {
	return uuid.New(), nil
}
func (pa ProjectAdapter) DeleteProject(ctx context.Context, projectUUID uuid.UUID) error {
	return nil
}
