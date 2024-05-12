package main

import (
	"github.com/aerosystems/customer-service/internal/config"
	"github.com/aerosystems/customer-service/internal/presenters/consumer"
	"github.com/sirupsen/logrus"
)

type App struct {
	log          *logrus.Logger
	cfg          *config.Config
	authConsumer *consumer.AuthSubscription
}

func NewApp(
	log *logrus.Logger,
	cfg *config.Config,
	authConsumer *consumer.AuthSubscription,
) *App {
	return &App{
		log:          log,
		cfg:          cfg,
		authConsumer: authConsumer,
	}
}
