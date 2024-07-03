//go:build wireinject
// +build wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	CustomErrors "github.com/aerosystems/customer-service/internal/common/custom_errors"
	"github.com/aerosystems/customer-service/internal/config"
	"github.com/aerosystems/customer-service/internal/infrastructure/adapters/broker"
	"github.com/aerosystems/customer-service/internal/infrastructure/repository/fire"
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
		wire.Bind(new(usecases.CustomerRepository), new(*fire.CustomerRepo)),
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
		ProvideBaseHandler,
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

func ProvideFireCustomerRepo(client *firestore.Client) *fire.CustomerRepo {
	panic(wire.Build(fire.NewCustomerRepo))
}

func ProvidePubSubClient(cfg *config.Config) *PubSub.Client {
	client, err := PubSub.NewClientWithAuth(cfg.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideHttpServer(log *logrus.Logger, cfg *config.Config, customErrorHandler *echo.HTTPErrorHandler, customerHandler *handlers.CustomerHandler) *HttpServer.Server {
	panic(wire.Build(HttpServer.NewServer))
}

func ProvideBaseHandler(log *logrus.Logger, cfg *config.Config) *handlers.BaseHandler {
	return handlers.NewBaseHandler(log, cfg.Mode)
}

func ProvideCustomerHandler(log *logrus.Logger, baseHandler *handlers.BaseHandler, customerUsecase handlers.CustomerUsecase) *handlers.CustomerHandler {
	panic(wire.Build(handlers.NewCustomerHandler))
}

func ProvideEchoErrorHandler(cfg *config.Config) *echo.HTTPErrorHandler {
	customErrorHandler := CustomErrors.NewEchoErrorHandler(cfg.Mode)
	return &customErrorHandler
}
