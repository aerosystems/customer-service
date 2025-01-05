package config

import (
	"github.com/spf13/viper"
)

const (
	defaultMode = "prod"
	defaultPort = 8080
)

type Config struct {
	Mode string
	Port int

	GcpProjectId                 string
	GoogleApplicationCredentials string

	ProjectServiceGRPCAddr      string
	SubscriptionServiceGRPCAddr string
}

func NewConfig() *Config {
	viper.AutomaticEnv()
	mode := viper.GetString("MODE")
	if mode == "" {
		mode = defaultMode
	}
	webPort := viper.GetInt("PORT")
	if webPort == 0 {
		webPort = defaultPort
	}
	return &Config{
		Mode:                         mode,
		Port:                         webPort,
		GcpProjectId:                 viper.GetString("GCP_PROJECT_ID"),
		GoogleApplicationCredentials: viper.GetString("GOOGLE_APPLICATION_CREDENTIALS"),
		ProjectServiceGRPCAddr:       viper.GetString("PRJCT_SERVICE_GRPC_ADDR"),
		SubscriptionServiceGRPCAddr:  viper.GetString("SBS_SERVICE_GRPC_ADDR"),
	}
}
