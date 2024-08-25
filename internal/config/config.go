package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/graytonio/flagops-config-storage/internal/facts"
	"github.com/graytonio/flagops-config-storage/internal/secrets"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func init() {
	viper.BindEnv("fact_provider")
	viper.BindEnv("secret_provider")

	viper.BindEnv("redis_uri")
}

func GetProviders() (facts.FactProvider, secrets.SecretProvider, error) {
	factProvider, err := createFactProvider()
	if err != nil {
	  return nil, nil, err
	}

	secretProvider, err := createSecretsProvider()
	if err != nil {
	  return nil, nil, err
	}

	return factProvider, secretProvider, nil
}

func createFactProvider() (facts.FactProvider, error) {
	switch viper.GetString("fact_provider") {
	case "redis":
		opts, err := redis.ParseURL(viper.GetString("redis_uri"))
		if err != nil {
		  return nil, err
		}

		return facts.NewRedisFactProvider(redis.NewClient(opts)), nil
	}

	return nil, fmt.Errorf("no such fact provider %s", viper.GetString("fact_provider"))
}

func createSecretsProvider() (secrets.SecretProvider, error) {
	switch viper.GetString("secret_provider") {
	case "asm":
		config, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
		  return nil, err
		}

		return secrets.NewASMSecretProvider(secretsmanager.NewFromConfig(config)), nil
	}

	return nil, fmt.Errorf("no such secret provider %s", viper.GetString("secret_provider"))
}