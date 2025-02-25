// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/aerosystems/common-service/clients/gcpclient"
	"github.com/aerosystems/common-service/logger"
	"github.com/aerosystems/common-service/presenters/httpserver"
	"github.com/aerosystems/customer-service/internal/adapters"
	"github.com/aerosystems/customer-service/internal/ports/http"
	"github.com/aerosystems/customer-service/internal/usecases"
	"github.com/sirupsen/logrus"
)

// Injectors from wire.go:

//go:generate wire
func InitApp() *App {
	logger := ProvideLogger()
	logrusLogger := ProvideLogrusLogger(logger)
	config := ProvideConfig()
	client := ProvideFirestoreClient(config)
	firestoreCustomerRepo := ProvideFirestoreCustomerRepo(client)
	subscriptionAdapter := ProvideSubscriptionAdapter(config)
	projectAdapter := ProvideProjectAdapter(config)
	checkmailAdapter := ProvideCheckmailAdapter(config)
	authClient := ProvideFirebaseAuthClient(config)
	firebaseAuthAdapter := ProvideFirebaseAuthAdapter(authClient)
	customerUsecase := ProvideCustomerUsecase(logrusLogger, firestoreCustomerRepo, subscriptionAdapter, projectAdapter, checkmailAdapter, firebaseAuthAdapter)
	handler := ProvideHandler(logrusLogger, customerUsecase)
	server := ProvideHTTPServer(config, logrusLogger, handler)
	app := ProvideApp(logrusLogger, config, server)
	return app
}

func ProvideApp(log *logrus.Logger, cfg *Config, httpServer *HTTPServer.Server) *App {
	app := NewApp(log, cfg, httpServer)
	return app
}

func ProvideLogger() *logger.Logger {
	loggerLogger := logger.NewLogger()
	return loggerLogger
}

func ProvideConfig() *Config {
	config := NewConfig()
	return config
}

func ProvideCustomerUsecase(log *logrus.Logger, customerRepo usecases.CustomerRepository, subscriptionAdapter usecases.SubscriptionAdapter, projectAdapter usecases.ProjectAdapter, checkmailAdapter usecases.CheckmailAdapter, firebaseAuthAdapter usecases.FirebaseAuthAdapter) *usecases.CustomerUsecase {
	customerUsecase := usecases.NewCustomerUsecase(log, customerRepo, subscriptionAdapter, projectAdapter, checkmailAdapter, firebaseAuthAdapter)
	return customerUsecase
}

func ProvideFirebaseAuthAdapter(client *auth.Client) *adapters.FirebaseAuthAdapter {
	firebaseAuthAdapter := adapters.NewFirebaseAuthAdapter(client)
	return firebaseAuthAdapter
}

func ProvideFirestoreCustomerRepo(client *firestore.Client) *adapters.FirestoreCustomerRepo {
	firestoreCustomerRepo := adapters.NewFirestoreCustomerRepo(client)
	return firestoreCustomerRepo
}

func ProvideHandler(log *logrus.Logger, customerUsecase HTTPServer.CustomerUsecase) *HTTPServer.Handler {
	handler := HTTPServer.NewHandler(customerUsecase)
	return handler
}

// wire.go:

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

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
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

func ProvideHTTPServer(cfg *Config, log *logrus.Logger, handler *HTTPServer.Handler) *HTTPServer.Server {
	return HTTPServer.NewHTTPServer(&HTTPServer.Config{
		Config: httpserver.Config{
			Host: cfg.Host,
			Port: cfg.Port,
		},
		Mode: cfg.Mode,
	}, log, handler)
}

func ProvideCheckmailAdapter(cfg *Config) *adapters.CheckmailAdapter {
	checkmailAdapter, err := adapters.NewCheckmailAdapter(cfg.CheckmailServiceGRPCAddr)
	if err != nil {
		panic(err)
	}
	return checkmailAdapter
}
