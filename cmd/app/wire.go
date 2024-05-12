//go:build wireinject
// +build wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/aerosystems/customer-service/internal/config"
	"github.com/aerosystems/customer-service/internal/infrastructure/adapters/rpc"
	"github.com/aerosystems/customer-service/internal/infrastructure/repository/fire"
	"github.com/aerosystems/customer-service/internal/presenters/consumer"
	"github.com/aerosystems/customer-service/internal/usecases"
	"github.com/aerosystems/customer-service/pkg/logger"
	PubSub "github.com/aerosystems/customer-service/pkg/pubsub"
	"github.com/aerosystems/customer-service/pkg/rpc_client"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(consumer.CustomerUsecase), new(*usecases.CustomerUsecase)),
		wire.Bind(new(usecases.CustomerRepository), new(*fire.CustomerRepo)),
		wire.Bind(new(usecases.SubsRepository), new(*RpcRepo.SubsRepo)),
		wire.Bind(new(usecases.ProjectRepository), new(*RpcRepo.ProjectRepo)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideLogrusLogger,
		ProvideFirestoreClient,
		ProvideCustomerUsecase,
		ProvideFireCustomerRepo,
		ProvideSubsRepo,
		ProvideProjectRepo,
		ProvideAuthConsumer,
		ProvidePubSubClient,
	))
}

func ProvideApp(log *logrus.Logger, cfg *config.Config, authConsumer *consumer.AuthSubscription) *App {
	panic(wire.Build(NewApp))
}

func ProvideLogger() *logger.Logger {
	panic(wire.Build(logger.NewLogger))
}

func ProvideConfig() *config.Config {
	panic(wire.Build(config.NewConfig))
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideCustomerUsecase(customerRepo usecases.CustomerRepository, projectRepo usecases.ProjectRepository, subsRepository usecases.SubsRepository) *usecases.CustomerUsecase {
	panic(wire.Build(usecases.NewCustomerUsecase))
}

func ProvideSubsRepo(cfg *config.Config) *RpcRepo.SubsRepo {
	rpcClient := RpcClient.NewClient("tcp", cfg.SubsServiceRPCAddress)
	return RpcRepo.NewSubsRepo(rpcClient)
}

func ProvideProjectRepo(cfg *config.Config) *RpcRepo.ProjectRepo {
	rpcClient := RpcClient.NewClient("tcp", cfg.ProjectServiceRpcAddress)
	return RpcRepo.NewProjectRepo(rpcClient)
}

func ProvideFirestoreClient(cfg *config.Config) *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GcpProjectId)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideFireCustomerRepo(client *firestore.Client) *fire.CustomerRepo {
	panic(wire.Build(fire.NewCustomerRepo))
}

func ProvideAuthConsumer(log *logrus.Logger, cfg *config.Config, client *PubSub.Client, customerUsecase consumer.CustomerUsecase) *consumer.AuthSubscription {
	return consumer.NewAuthSubscription(log, client, cfg.AuthTopicId, cfg.AuthSubName, customerUsecase)
}

func ProvidePubSubClient(cfg *config.Config) *PubSub.Client {
	var client *PubSub.Client
	var err error
	switch cfg.Mode {
	case "dev":
		client, err = PubSub.NewClient(cfg.GcpProjectId)
	default:
		client, err = PubSub.NewClientWithAuth(cfg.GoogleApplicationCredentials)
	}
	if err != nil {
		panic(err)
	}
	return client
}
