//go:build wireinject
// +build wireinject

package app

import (
	"firebase.google.com/go/v4/auth"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/aerosystems/common-service/clients/gcpclient"
	"github.com/aerosystems/common-service/clients/gormclient"
	"github.com/aerosystems/common-service/logger"
	"github.com/aerosystems/customer-service/internal/adapters"
	HTTPServer "github.com/aerosystems/customer-service/internal/ports/http"
	"github.com/aerosystems/customer-service/internal/usecases"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(HTTPServer.CustomerUsecase), new(*usecases.CustomerUsecase)),
		wire.Bind(new(usecases.CustomerRepository), new(*adapters.CustomerPostgresRepo)),
		wire.Bind(new(usecases.SubscriptionAdapter), new(*adapters.SubscriptionAdapter)),
		wire.Bind(new(usecases.ProjectAdapter), new(*adapters.ProjectAdapter)),
		wire.Bind(new(usecases.FirebaseAuthAdapter), new(*adapters.FirebaseAuthAdapter)),
		wire.Bind(new(usecases.CheckmailAdapter), new(*adapters.CheckmailAdapter)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideLogrusLogger,
		ProvideGORMPostgres,
		ProvideCustomerUsecase,
		ProvideCustomerPostgresRepo,
		ProvideHTTPServer,
		ProvideHandler,
		ProvideSubscriptionAdapter,
		ProvideCheckmailAdapter,
		ProvideProjectAdapter,
		ProvideFirebaseAuthClient,
		ProvideFirebaseAuthAdapter,
	))
}

//go:generate wire
func InitAppMigration() *AppMigration {
	panic(wire.Build(
		ProvideAppMigration,
		ProvideMigration,
		ProvideLogger,
		ProvideConfig,
		ProvideLogrusLogger,
		ProvideGORMPostgres,
	))
}

func ProvideApp(log *logrus.Logger, cfg *Config, httpServer *HTTPServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideAppMigration(log *logrus.Logger, cfg *Config, migration *adapters.Migration) *AppMigration {
	panic(wire.Build(NewAppMigration))
}

func ProvideSubscriptionAdapter(cfg *Config) *adapters.SubscriptionAdapter {
	subscriptionAdapter, err := adapters.NewSubscriptionAdapter(&cfg.SubscriptionServiceGRPC)
	if err != nil {
		panic(err)
	}
	return subscriptionAdapter
}

func ProvideMigration(db *gorm.DB) *adapters.Migration {
	panic(wire.Build(adapters.NewMigration))
}

func ProvideProjectAdapter(cfg *Config) *adapters.ProjectAdapter {
	projectAdapter, err := adapters.NewProjectAdapter(&cfg.ProjectServiceGRPC)
	if err != nil {
		panic(err)
	}
	return projectAdapter
}

func ProvideCheckmailAdapter(cfg *Config) *adapters.CheckmailAdapter {
	checkmailAdapter, err := adapters.NewCheckmailAdapter(&cfg.CheckmailServiceGRPC)
	if err != nil {
		panic(err)
	}
	return checkmailAdapter
}

func ProvideLogger() *logger.Logger {
	panic(wire.Build(logger.NewLogger))
}

func ProvideConfig() *Config {
	panic(wire.Build(NewConfig))
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideGORMPostgres(log *logrus.Logger, cfg *Config) *gorm.DB {
	db := gormclient.NewPostgresDB(log, &cfg.Postgres)
	return db
}

func ProvideCustomerUsecase(log *logrus.Logger, customerRepo usecases.CustomerRepository, subscriptionAdapter usecases.SubscriptionAdapter, projectAdapter usecases.ProjectAdapter, checkmailAdapter usecases.CheckmailAdapter, firebaseAuthAdapter usecases.FirebaseAuthAdapter) *usecases.CustomerUsecase {
	panic(wire.Build(usecases.NewCustomerUsecase))
}

func ProvideFirebaseAuthAdapter(client *auth.Client) *adapters.FirebaseAuthAdapter {
	panic(wire.Build(adapters.NewFirebaseAuthAdapter))
}

func ProvideFirebaseAuthClient(cfg *Config) *auth.Client {
	client, err := gcpclient.NewFirebaseClient(&cfg.Firebase)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideCustomerPostgresRepo(db *gorm.DB) *adapters.CustomerPostgresRepo {
	panic(wire.Build(adapters.NewCustomerPostgresRepo))
}

func ProvideHTTPServer(cfg *Config, log *logrus.Logger, handler *HTTPServer.Handler) *HTTPServer.Server {
	return HTTPServer.NewHTTPServer(&cfg.HTTPServer, cfg.Debug, log, handler)
}

func ProvideHandler(log *logrus.Logger, customerUsecase HTTPServer.CustomerUsecase) *HTTPServer.Handler {
	panic(wire.Build(HTTPServer.NewHandler))
}
