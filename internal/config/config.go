package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mode                         string
	WebPort                      int
	GcpProjectId                 string
	GoogleApplicationCredentials string
	ProjectServiceRpcAddress     string
	SubsServiceRPCAddress        string
	AuthTopicId                  string
	AuthSubName                  string
}

func NewConfig() *Config {
	viper.AutomaticEnv()
	return &Config{
		Mode:                         viper.GetString("CSTMR_MODE"),
		WebPort:                      viper.GetInt("PORT"),
		GcpProjectId:                 viper.GetString("GCP_PROJECT_ID"),
		GoogleApplicationCredentials: viper.GetString("GOOGLE_APPLICATION_CREDENTIALS"),
		ProjectServiceRpcAddress:     viper.GetString("CSTMR_PROJECT_SERVICE_RPC_ADDR"),
		SubsServiceRPCAddress:        viper.GetString("CSTMR_SUBS_SERVICE_RPC_ADDR"),
		AuthTopicId:                  viper.GetString("CSTMR_AUTH_TOPIC_ID"),
		AuthSubName:                  viper.GetString("CSTMR_AUTH_SUB_NAME"),
	}
}
