//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"

	"github.com/aerosystems/common-service/clients/gcpclient"
	"github.com/aerosystems/common-service/logger"
	"github.com/aerosystems/common-service/presenters/httpserver"
	"github.com/aerosystems/customer-service/internal/adapters"
	HTTPServer "github.com/aerosystems/customer-service/internal/ports/http"
	"github.com/aerosystems/customer-service/internal/usecases"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(HTTPServer.CustomerUsecase), new(*usecases.CustomerUsecase)),
		wire.Bind(new(usecases.CustomerRepository), new(*adapters.FirestoreCustomerRepo)),
		wire.Bind(new(usecases.SubscriptionAdapter), new(*adapters.SubscriptionAdapter)),
		wire.Bind(new(usecases.ProjectAdapter), new(*adapters.ProjectAdapter)),
		wire.Bind(new(usecases.FirebaseAuthAdapter), new(*adapters.FirebaseAuthAdapter)),
		wire.Bind(new(usecases.CheckmailAdapter), new(*adapters.CheckmailAdapter)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideLogrusLogger,
		ProvideFirestoreClient,
		ProvideCustomerUsecase,
		ProvideFirestoreCustomerRepo,
		ProvideHTTPServer,
		ProvideHandler,
		ProvideSubscriptionAdapter,
		ProvideProjectAdapter,
		ProvideFirebaseAuthClient,
		ProvideFirebaseAuthAdapter,
		ProvideCheckmailAdapter,
	))
}

func ProvideApp(log *logrus.Logger, cfg *Config, httpServer *HTTPServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideSubscriptionAdapter(cfg *Config) *adapters.SubscriptionAdapter {
	subscriptionAdapter, err := adapters.NewSubscriptionAdapter(cfg.SubscriptionServiceGRPCAddr)
	if err != nil {
		panic(err)
	}
	return subscriptionAdapter
}

func ProvideProjectAdapter(cfg *Config) *adapters.ProjectAdapter {
	projectAdapter, err := adapters.NewProjectAdapter(cfg.ProjectServiceGRPCAddr)
	if err != nil {
		panic(err)
	}
	return projectAdapter
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

func ProvideCustomerUsecase(log *logrus.Logger, customerRepo usecases.CustomerRepository, subscriptionAdapter usecases.SubscriptionAdapter, projectAdapter usecases.ProjectAdapter, checkmailAdapter usecases.CheckmailAdapter, firebaseAuthAdapter usecases.FirebaseAuthAdapter) *usecases.CustomerUsecase {
	panic(wire.Build(usecases.NewCustomerUsecase))
}

func ProvideFirebaseAuthAdapter(client *auth.Client) *adapters.FirebaseAuthAdapter {
	panic(wire.Build(adapters.NewFirebaseAuthAdapter))
}

func ProvideFirebaseAuthClient(cfg *Config) *auth.Client {
	client, err := gcpclient.NewFirebaseClient(cfg.GcpProjectId, cfg.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideFirestoreClient(cfg *Config) *firestore.Client {
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

func ProvideHTTPServer(cfg *Config, log *logrus.Logger, handler *HTTPServer.Handler) *HTTPServer.Server {
	return HTTPServer.NewHTTPServer(&HTTPServer.Config{
		Config: httpserver.Config{
			Host: cfg.Host,
			Port: cfg.Port,
		},
		Mode: cfg.Mode,
	}, log, handler)
}

func ProvideHandler(log *logrus.Logger, customerUsecase HTTPServer.CustomerUsecase) *HTTPServer.Handler {
	panic(wire.Build(HTTPServer.NewHandler))
}

func ProvideCheckmailAdapter(cfg *Config) *adapters.CheckmailAdapter {
	checkmailAdapter, err := adapters.NewCheckmailAdapter(cfg.CheckmailServiceGRPCAddr)
	if err != nil {
		panic(err)
	}
	return checkmailAdapter
}
