package config

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
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
	SECRET_NAME           string `mapstructure:"SECRET_NAME"`
}

func LoadConfig() (cfg Config, err error) {
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	paths := []string{".env", "../.env", "/app/.env"}
	loaded := false

	for _, path := range paths {
		viper.SetConfigFile(path)
		if err := viper.ReadInConfig(); err == nil {
			log.Printf("Loaded configuration from %s", path)
			loaded = true
			break
		}
	}

	if loaded {
		err = viper.Unmarshal(&cfg)
		return cfg, err
	}

	log.Println("Falling back to AWS Secrets Manager for configuration")
	secretName := os.Getenv("SECRET_NAME")
	if secretName == "" {
		secretName = "zyra/prod/api-gateway/env"
	}

	err = loadFromSecretsManager(&cfg, secretName)
	return cfg, err
}

func loadFromSecretsManager(cfg *Config, secretName string) error {
	ctx := context.TODO()

	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := secretsmanager.NewFromConfig(awsCfg)

	result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(*result.SecretString), cfg)
}
