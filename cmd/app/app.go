package main

import (
	"github.com/aerosystems/customer-service/internal/config"
	HttpServer "github.com/aerosystems/customer-service/internal/presenters/http"
	"github.com/sirupsen/logrus"
)

type App struct {
	log        *logrus.Logger
	cfg        *config.Config
	httpServer *HttpServer.Server
}

func NewApp(
	log *logrus.Logger,
	cfg *config.Config,
	httpServer *HttpServer.Server,
) *App {
	return &App{
		log:        log,
		cfg:        cfg,
		httpServer: httpServer,
	}
}
