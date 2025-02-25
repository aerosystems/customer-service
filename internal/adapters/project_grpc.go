package adapters

import (
	"context"
	"github.com/aerosystems/common-service/clients/grpcclient"
	"github.com/aerosystems/common-service/gen/protobuf/project"
	"github.com/google/uuid"
)

type ProjectAdapter struct {
	client project.ProjectServiceClient
}

func NewProjectAdapter(address string) (*ProjectAdapter, error) {
	conn, err := grpcclient.NewGRPCConn(address)
	if err != nil {
		return nil, err
	}
	return &ProjectAdapter{
		client: project.NewProjectServiceClient(conn),
	}, nil
}

func (pa ProjectAdapter) CreateDefaultProject(ctx context.Context, customerUUID uuid.UUID) (uuid.UUID, string, error) {
	resp, err := pa.client.CreateDefaultProject(ctx, &project.CreateDefaultProjectRequest{
		CustomerUuid: customerUUID.String(),
	})
	if err != nil {
		return uuid.Nil, "", err
	}
	projectUuid, err := uuid.Parse(resp.ProjectUuid)
	if err != nil {
		return uuid.Nil, "", err
	}
	return projectUuid, resp.ProjectToken, nil
}

func (pa ProjectAdapter) DeleteProject(ctx context.Context, projectUUID uuid.UUID) error {
	_, err := pa.client.DeleteProject(ctx, &project.DeleteProjectRequest{
		ProjectUuid: projectUUID.String(),
	})
	return err
}
