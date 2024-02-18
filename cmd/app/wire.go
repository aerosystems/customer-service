//go:build wireinject
// +build wireinject

package main

import (
	"github.com/aerosystems/customer-service/internal/config"
	"github.com/aerosystems/customer-service/internal/http"
	"github.com/aerosystems/customer-service/internal/infrastructure/rest"
	"github.com/aerosystems/customer-service/internal/infrastructure/rpc"
	"github.com/aerosystems/customer-service/internal/models"
	"github.com/aerosystems/customer-service/internal/repository/pg"
	"github.com/aerosystems/customer-service/internal/repository/rpc"
	"github.com/aerosystems/customer-service/internal/usecases"
	"github.com/aerosystems/customer-service/pkg/gorm_postgres"
	"github.com/aerosystems/customer-service/pkg/logger"
	"github.com/aerosystems/customer-service/pkg/oauth"
	"github.com/aerosystems/customer-service/pkg/rpc_client"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(rest.CustomerUsecase), new(*usecases.CustomerUsecase)),
		wire.Bind(new(RpcServer.CustomerUsecase), new(*usecases.CustomerUsecase)),
		wire.Bind(new(usecases.CustomerRepository), new(*pg.CustomerRepo)),
		wire.Bind(new(usecases.SubsRepository), new(*RpcRepo.SubsRepo)),
		wire.Bind(new(usecases.ProjectRepository), new(*RpcRepo.ProjectRepo)),
		wire.Bind(new(HttpServer.TokenService), new(*OAuthService.AccessTokenService)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideHttpServer,
		ProvideRpcServer,
		ProvideLogrusLogger,
		ProvideLogrusEntry,
		ProvideGormPostgres,
		ProvideBaseHandler,
		ProvideCustomerHandler,
		ProvideCustomerUsecase,
		ProvideCustomerRepo,
		ProvideSubsRepo,
		ProvideProjectRepo,
		ProvideAccessTokenService,
	))
}

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HttpServer.Server, rpcServer *RpcServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideLogger() *logger.Logger {
	panic(wire.Build(logger.NewLogger))
}

func ProvideConfig() *config.Config {
	panic(wire.Build(config.NewConfig))
}

func ProvideHttpServer(log *logrus.Logger, cfg *config.Config, customerHandler *rest.CustomerHandler, tokenService HttpServer.TokenService) *HttpServer.Server {
	panic(wire.Build(HttpServer.NewServer))
}

func ProvideRpcServer(log *logrus.Logger, customerUsecase RpcServer.CustomerUsecase) *RpcServer.Server {
	panic(wire.Build(RpcServer.NewServer))
}

func ProvideLogrusEntry(log *logger.Logger) *logrus.Entry {
	return logrus.NewEntry(log.Logger)
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideGormPostgres(e *logrus.Entry, cfg *config.Config) *gorm.DB {
	db := GormPostgres.NewClient(e, cfg.PostgresDSN)
	if err := db.AutoMigrate(&models.Customer{}); err != nil { // TODO: Move to migration
		panic(err)
	}
	return db
}

func ProvideBaseHandler(log *logrus.Logger, cfg *config.Config) *rest.BaseHandler {
	return rest.NewBaseHandler(log, cfg.Mode)
}

func ProvideCustomerHandler(baseHandler *rest.BaseHandler, customerUsecase rest.CustomerUsecase) *rest.CustomerHandler {
	panic(wire.Build(rest.NewCustomerHandler))
}

func ProvideCustomerUsecase(customerRepo usecases.CustomerRepository, projectRepo usecases.ProjectRepository, subsRepository usecases.SubsRepository) *usecases.CustomerUsecase {
	panic(wire.Build(usecases.NewCustomerUsecase))
}

func ProvideCustomerRepo(db *gorm.DB) *pg.CustomerRepo {
	panic(wire.Build(pg.NewCustomerRepo))
}

func ProvideSubsRepo(cfg *config.Config) *RpcRepo.SubsRepo {
	rpcClient := RPCClient.NewClient("tcp", cfg.SubsServiceRPCAddress)
	return RpcRepo.NewSubsRepo(rpcClient)
}

func ProvideProjectRepo(cfg *config.Config) *RpcRepo.ProjectRepo {
	rpcClient := RPCClient.NewClient("tcp", cfg.ProjectServiceRPCAddress)
	return RpcRepo.NewProjectRepo(rpcClient)
}

func ProvideAccessTokenService(cfg *config.Config) *OAuthService.AccessTokenService {
	return OAuthService.NewAccessTokenService(cfg.AccessSecret)
}
