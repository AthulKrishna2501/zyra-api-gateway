package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port                string `mapstructure:"PORT"`
	AUTH_SVC_URL        string `mapstructure:"AUTH_SVC_URL"`
	CLIENT_SVC_URL      string `mapstructure:"CLIENT_SVC_URL"`
	ADMIN_SVC_URL       string `mapstructure:"ADMIN_SVC_URL"`
	VENDOR_SVC_URL      string `mapstructure:"VENDOR_SVC_URL"`
	MESSAGE_SERVICE_URL string `mapstructure:"MESSAGE_SVC_URL"`
	CHAT_SERVICE_URL    string `mapstructure:"CHAT_SVC_URL"`
	RABBITMQ_URL        string `mapstructure:"RABBITMQ_URL"`
}

func LoadConfig() (cfg Config, err error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.SetConfigFile("../.env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	viper.AutomaticEnv()

	err = viper.Unmarshal(&cfg)

	return
}
