package app

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/aerosystems/common-service/clients/gcpclient"
	"github.com/aerosystems/common-service/clients/gormclient"
	"github.com/aerosystems/common-service/clients/grpcclient"
	"github.com/aerosystems/common-service/presenters/httpserver"
)

type Config struct {
	Debug                   bool
	HTTPServer              httpserver.Config
	Postgres                gormclient.Config
	ProjectServiceGRPC      grpcclient.Config
	SubscriptionServiceGRPC grpcclient.Config
	CheckmailServiceGRPC    grpcclient.Config
	Firebase                gcpclient.FirebaseConfig
}

func NewConfig() *Config {
	viper.Reset()
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("failed to unmarshal cfg: %w", err))
	}

	return &cfg
}
