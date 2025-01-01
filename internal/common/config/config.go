package config

import (
	"github.com/spf13/viper"
)

const (
	defaultMode    = "prod"
	defaultWebPort = 8080
)

type Config struct {
	Mode                                string
	WebPort                             int
	GcpProjectId                        string
	GoogleApplicationCredentials        string
	SubscriptionTopicId                 string
	SubscriptionSubName                 string
	SubscriptionCreateFreeTrialEndpoint string
	SubscriptionServiceApiKey           string
}

func NewConfig() *Config {
	viper.AutomaticEnv()
	mode := viper.GetString("CSTMR_MODE")
	if mode == "" {
		mode = defaultMode
	}
	webPort := viper.GetInt("PORT")
	if webPort == 0 {
		webPort = defaultWebPort
	}
	return &Config{
		Mode:                                mode,
		WebPort:                             webPort,
		GcpProjectId:                        viper.GetString("GCP_PROJECT_ID"),
		GoogleApplicationCredentials:        viper.GetString("GOOGLE_APPLICATION_CREDENTIALS"),
		SubscriptionTopicId:                 viper.GetString("CSTMR_SUBSCRIPTION_TOPIC_ID"),
		SubscriptionSubName:                 viper.GetString("CSTMR_SUBSCRIPTION_SUB_NAME"),
		SubscriptionCreateFreeTrialEndpoint: viper.GetString("CSTMR_SUBSCRIPTION_CREATE_FREE_TRIAL_ENDPOINT"),
		SubscriptionServiceApiKey:           viper.GetString("SBS_API_KEY"),
	}
}
