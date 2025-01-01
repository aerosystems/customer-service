//go:build wireinject
// +build wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/aerosystems/customer-service/internal/common/config"
	CustomErrors "github.com/aerosystems/customer-service/internal/common/custom_errors"
	"github.com/aerosystems/customer-service/internal/adapters/broker"
	FirestoreRepo "github.com/aerosystems/customer-service/internal/adapters/firestore_repo"
	HttpServer "github.com/aerosystems/customer-service/internal/presenters/http"
	"github.com/aerosystems/customer-service/internal/presenters/http/handlers"
	"github.com/aerosystems/customer-service/internal/usecases"
	"github.com/aerosystems/customer-service/pkg/logger"
	PubSub "github.com/aerosystems/customer-service/pkg/pubsub"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(handlers.CustomerUsecase), new(*usecases.CustomerUsecase)),
		wire.Bind(new(usecases.CustomerRepository), new(*FirestoreRepo.CustomerRepo)),
		wire.Bind(new(usecases.SubscriptionEventsAdapter), new(*broker.SubscriptionEventsAdapter)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideLogrusLogger,
		ProvideFirestoreClient,
		ProvideCustomerUsecase,
		ProvideFireCustomerRepo,
		ProvideHttpServer,
		ProvideCustomerHandler,
		ProvidePubSubClient,
		ProvideSubscriptionEventsAdapter,
		ProvideEchoErrorHandler,
	))
}

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HttpServer.Server) *App {
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

func ProvideCustomerUsecase(log *logrus.Logger, customerRepo usecases.CustomerRepository, subscriptionEventsAdapter usecases.SubscriptionEventsAdapter) *usecases.CustomerUsecase {
	panic(wire.Build(usecases.NewCustomerUsecase))
}

func ProvideSubscriptionEventsAdapter(pubSubClient *PubSub.Client, cfg *config.Config) *broker.SubscriptionEventsAdapter {
	return broker.NewSubscriptionEventsAdapter(pubSubClient, cfg.SubscriptionTopicId, cfg.SubscriptionSubName, cfg.SubscriptionCreateFreeTrialEndpoint, cfg.SubscriptionServiceApiKey)
}

func ProvideFirestoreClient(cfg *config.Config) *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GcpProjectId)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideFireCustomerRepo(client *firestore.Client) *FirestoreRepo.CustomerRepo {
	panic(wire.Build(FirestoreRepo.NewCustomerRepo))
}

func ProvidePubSubClient(cfg *config.Config) *PubSub.Client {
	client, err := PubSub.NewClientWithAuth(cfg.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideHttpServer(cfg *config.Config, log *logrus.Logger, customErrorHandler *echo.HTTPErrorHandler, customerHandler *handlers.FirebaseHandler) *HttpServer.Server {
	return HttpServer.NewServer(cfg.WebPort, log, customErrorHandler, customerHandler)
}

func ProvideCustomerHandler(log *logrus.Logger, customerUsecase handlers.CustomerUsecase) *handlers.FirebaseHandler {
	panic(wire.Build(handlers.NewFirebaseHandler))
}

func ProvideEchoErrorHandler(cfg *config.Config) *echo.HTTPErrorHandler {
	customErrorHandler := CustomErrors.NewEchoErrorHandler(cfg.Mode)
	return &customErrorHandler
}
