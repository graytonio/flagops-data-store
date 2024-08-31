package config

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/graytonio/flagops-data-storage/internal/facts"
	"github.com/graytonio/flagops-data-storage/internal/secrets"
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
		cfgOpts := []func(*config.LoadOptions) error{}
		if os.Getenv("LOCAL_AWS_ENDPOINT") != "" {
			cfgOpts = append(cfgOpts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("accessKey", "secretKey", "token")))
		}

		config, err := config.LoadDefaultConfig(context.Background(), cfgOpts...)
		if err != nil {
		  return nil, err
		}

		return secrets.NewASMSecretProvider(secretsmanager.NewFromConfig(config, func(o *secretsmanager.Options) {
			if os.Getenv("LOCAL_AWS_ENDPOINT") != "" {
				o.BaseEndpoint = aws.String(os.Getenv("LOCAL_AWS_ENDPOINT"))
			}
		})), nil
	}

	return nil, fmt.Errorf("no such secret provider %s", viper.GetString("secret_provider"))
}