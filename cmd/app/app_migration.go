package app

import (
	"github.com/sirupsen/logrus"

	"github.com/aerosystems/customer-service/internal/adapters"
)

type Migration struct {
	log       *logrus.Logger
	cfg       *Config
	migration *adapters.Migration
}

func NewAppMigration(
	log *logrus.Logger,
	cfg *Config,
	migration *adapters.Migration,
) *Migration {
	return &Migration{
		log:       log,
		cfg:       cfg,
		migration: migration,
	}
}
