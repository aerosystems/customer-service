//go:build wireinject
// +build wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/aerosystems/common-service/pkg/gcp"
	"github.com/aerosystems/common-service/pkg/logger"
	"github.com/aerosystems/customer-service/internal/adapters"
	"github.com/aerosystems/customer-service/internal/common/config"
	CustomErrors "github.com/aerosystems/customer-service/internal/common/custom_errors"
	HTTPServer "github.com/aerosystems/customer-service/internal/presenters/http"
	"github.com/aerosystems/customer-service/internal/usecases"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(HTTPServer.CustomerUsecase), new(*usecases.CustomerUsecase)),
		wire.Bind(new(usecases.CustomerRepository), new(*adapters.FirestoreCustomerRepo)),
		wire.Bind(new(usecases.SubscriptionAdapter), new(*adapters.SubscriptionAdapter)),
		wire.Bind(new(usecases.ProjectAdapter), new(*adapters.ProjectAdapter)),
		wire.Bind(new(usecases.FirebaseAuthAdapter), new(*adapters.FirebaseAuthAdapter)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideLogrusLogger,
		ProvideFirestoreClient,
		ProvideCustomerUsecase,
		ProvideFirestoreCustomerRepo,
		ProvideHttpServer,
		ProvideCustomerHandler,
		ProvideEchoErrorHandler,
		ProvideSubscriptionAdapter,
		ProvideProjectAdapter,
		ProvideFirebaseAuthClient,
		ProvideFirebaseAuthAdapter,
	))
}

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HTTPServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideSubscriptionAdapter(cfg *config.Config) *adapters.SubscriptionAdapter {
	subscriptionAdapter, err := adapters.NewSubscriptionAdapter(cfg.SubscriptionServiceGRPCAddr)
	if err != nil {
		panic(err)
	}
	return subscriptionAdapter
}

func ProvideProjectAdapter(cfg *config.Config) *adapters.ProjectAdapter {
	projectAdapter, err := adapters.NewProjectAdapter(cfg.ProjectServiceGRPCAddr)
	if err != nil {
		panic(err)
	}
	return projectAdapter
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

func ProvideCustomerUsecase(log *logrus.Logger, customerRepo usecases.CustomerRepository, subscriptionAdapter usecases.SubscriptionAdapter, projectAdapter usecases.ProjectAdapter, firebaseAuthAdapter usecases.FirebaseAuthAdapter) *usecases.CustomerUsecase {
	panic(wire.Build(usecases.NewCustomerUsecase))
}

func ProvideFirebaseAuthAdapter(client *auth.Client) *adapters.FirebaseAuthAdapter {
	panic(wire.Build(adapters.NewFirebaseAuthAdapter))
}

func ProvideFirebaseAuthClient(cfg *config.Config) *auth.Client {
	client, err := gcp.NewFirebaseClient(cfg.GcpProjectId, cfg.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideFirestoreClient(cfg *config.Config) *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GcpProjectId)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideFirestoreCustomerRepo(client *firestore.Client) *adapters.FirestoreCustomerRepo {
	panic(wire.Build(adapters.NewFirestoreCustomerRepo))
}

func ProvideHttpServer(cfg *config.Config, log *logrus.Logger, customErrorHandler *echo.HTTPErrorHandler, customerHandler *HTTPServer.FirebaseHandler) *HTTPServer.Server {
	return HTTPServer.NewServer(cfg.Port, log, customErrorHandler, customerHandler)
}

func ProvideCustomerHandler(log *logrus.Logger, customerUsecase HTTPServer.CustomerUsecase) *HTTPServer.FirebaseHandler {
	panic(wire.Build(HTTPServer.NewFirebaseHandler))
}

func ProvideEchoErrorHandler(cfg *config.Config) *echo.HTTPErrorHandler {
	customErrorHandler := CustomErrors.NewEchoErrorHandler(cfg.Mode)
	return &customErrorHandler
}
