//go:build wireinject
// +build wireinject

package main

import (
	"github.com/aerosystems/customer-service/internal/config"
	"github.com/aerosystems/customer-service/internal/http"
	"github.com/aerosystems/customer-service/internal/infrastructure/rest"
	"github.com/aerosystems/customer-service/internal/infrastructure/rpc"
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
		wire.Bind(new(RPCServer.CustomerUsecase), new(*usecases.CustomerUsecase)),
		wire.Bind(new(usecases.CustomerRepository), new(*pg.CustomerRepo)),
		wire.Bind(new(usecases.SubsRepository), new(*rpcRepo.SubsRepo)),
		wire.Bind(new(usecases.ProjectRepository), new(*rpcRepo.ProjectRepo)),
		wire.Bind(new(HTTPServer.TokenService), new(*OAuthService.AccessTokenService)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideHTTPServer,
		ProvideRPCServer,
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

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HTTPServer.Server, rpcServer *RPCServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideLogger() *logger.Logger {
	panic(wire.Build(logger.NewLogger))
}

func ProvideConfig() *config.Config {
	panic(wire.Build(config.NewConfig))
}

func ProvideHTTPServer(log *logrus.Logger, cfg *config.Config, customerHandler *rest.CustomerHandler, tokenService HTTPServer.TokenService) *HTTPServer.Server {
	panic(wire.Build(HTTPServer.NewServer))
}

func ProvideRPCServer(log *logrus.Logger, customerUsecase RPCServer.CustomerUsecase) *RPCServer.Server {
	panic(wire.Build(RPCServer.NewServer))
}

func ProvideLogrusEntry(log *logger.Logger) *logrus.Entry {
	return logrus.NewEntry(log.Logger)
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideGormPostgres(e *logrus.Entry, cfg *config.Config) *gorm.DB {
	return GormPostgres.NewClient(e, cfg.PostgresDSN)
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

func ProvideSubsRepo(cfg *config.Config) *rpcRepo.SubsRepo {
	rpcClient := RPCClient.NewClient("tcp", cfg.SubsServiceRPCAddress)
	return rpcRepo.NewSubsRepo(rpcClient)
}

func ProvideProjectRepo(cfg *config.Config) *rpcRepo.ProjectRepo {
	rpcClient := RPCClient.NewClient("tcp", cfg.ProjectServiceRPCAddress)
	return rpcRepo.NewProjectRepo(rpcClient)
}

func ProvideAccessTokenService(cfg *config.Config) *OAuthService.AccessTokenService {
	return OAuthService.NewAccessTokenService(cfg.AccessSecret)
}
