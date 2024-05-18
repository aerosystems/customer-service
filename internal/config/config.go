package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mode                         string
	WebPort                      int
	GcpProjectId                 string
	GoogleApplicationCredentials string
	SubscriptionTopicId          string
}

func NewConfig() *Config {
	viper.AutomaticEnv()
	return &Config{
		Mode:                         viper.GetString("CSTMR_MODE"),
		WebPort:                      viper.GetInt("PORT"),
		GcpProjectId:                 viper.GetString("GCP_PROJECT_ID"),
		GoogleApplicationCredentials: viper.GetString("GOOGLE_APPLICATION_CREDENTIALS"),
		SubscriptionTopicId:          viper.GetString("CSTMR_SUBSCRIPTION_TOPIC_ID"),
	}
}
