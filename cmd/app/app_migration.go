package app

import (
	"github.com/sirupsen/logrus"

	"github.com/aerosystems/customer-service/internal/adapters"
)

type AppMigration struct {
	log       *logrus.Logger
	cfg       *Config
	migration *adapters.Migration
}

func NewAppMigration(
	log *logrus.Logger,
	cfg *Config,
	migration *adapters.Migration,
) *AppMigration {
	return &AppMigration{
		log:       log,
		cfg:       cfg,
		migration: migration,
	}
}
