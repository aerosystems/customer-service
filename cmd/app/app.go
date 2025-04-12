package main

import (
	"github.com/sirupsen/logrus"

	HTTPServer "github.com/aerosystems/customer-service/internal/ports/http"
)

type App struct {
	log        *logrus.Logger
	cfg        *Config
	httpServer *HTTPServer.Server
}

func NewApp(
	log *logrus.Logger,
	cfg *Config,
	httpServer *HTTPServer.Server,
) *App {
	return &App{
		log:        log,
		cfg:        cfg,
		httpServer: httpServer,
	}
}
