package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port                  string `mapstructure:"PORT"`
	AUTH_SVC_URL          string `mapstructure:"AUTH_SVC_URL"`
	CLIENT_SVC_URL        string `mapstructure:"CLIENT_SVC_URL"`
	ADMIN_SVC_URL         string `mapstructure:"ADMIN_SVC_URL"`
	VENDOR_SVC_URL        string `mapstructure:"VENDOR_SVC_URL"`
	CHAT_SERVICE_URL      string `mapstructure:"CHAT_SVC_URL"`
	STRIPE_SECRET_KEY     string `mapstructure:"STRIPE_SECRET_KEY"`
	STRIPE_WEBHOOK_SECRET string `mapstructure:"STRIPE_WEBHOOK_SECRET"`
	ADMIN_EMAIL           string `mapstructure:"ADMIN_EMAIL"`
}

func LoadConfig() (cfg Config, err error) {
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Loaded .env from the current directory")
	} else {
		log.Println("Could not load .env from current directory, trying parent...")

		viper.SetConfigFile("../.env")
		if err := viper.ReadInConfig(); err == nil {
			log.Println("Loaded .env from parent directory")
		} else {
			log.Println("Could not load .env from parent directory, trying absolute path...")

			viper.SetConfigFile("/app/.env")
			if err := viper.ReadInConfig(); err == nil {
				log.Println("Loaded .env from absolute path (/app/.env)")
			} else {
				log.Fatalf("Error loading .env file: %v", err)
			}
		}
	}

	err = viper.Unmarshal(&cfg)
	return
}
