package config

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	FactsProviderOptions FactsProviderOptions `mapstructure:"facts"`
	SecretsProviderOptions SecretsProviderOptions `mapstructure:"secrets"`
	
	UserDatabaseOptions UserDatabaseOptions `mapstructure:"user_db"`
	OAuthOptions OAuthOptions `mapstructure:"oauth"`
}

type FactsProviderOptions struct {
	Provider string `mapstructure:"provider"`
	RedisURI string `mapstructure:"redis_uri"`
}

type SecretsProviderOptions struct {
	Provider string `mapstructure:"provider"`
	ASMDeletionRecoveryDays int `mapstructure:"asm_deletion_recovery_days"`
}

type UserDatabaseOptions struct {
	PostgresDSN string `mapstructure:"dsn"`
	JWTSecret string `mapstructure:"signing_secret"`
	AccessTokenExpirationMinutes int `mapstructure:"access_token_expiration_minutes"`
	RefreshTokenExpirationMinutes int `mapstructure:"refresh_token_expiration_minutes"`
}

type OAuthOptions struct {
	Provider string `mapstructure:"provider"`
	Hostname string `mapstructure:"hostname"`

	GithubClientKey string `mapstructure:"github_client_key"`
	GithubClientSecret string `mapstructure:"github_client_secret"`
}

func ParseConfig() (*Config, error) {
	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Set default values here
	conf := Config{
		FactsProviderOptions: FactsProviderOptions{
			Provider: "redis",
		},
		SecretsProviderOptions: SecretsProviderOptions{
			Provider: "asm",
			ASMDeletionRecoveryDays: 7,
		},
		UserDatabaseOptions: UserDatabaseOptions{
			JWTSecret: "flagops-salt",
			AccessTokenExpirationMinutes: 15,
			RefreshTokenExpirationMinutes: 720,
		},
	}

	err := v.Unmarshal(&conf)
	if err != nil {
	  return nil, err
	}

	logrus.Infof("%+v", conf)

	return &conf, nil
}

// func GetProviders(conf Config) (facts.FactProvider, secrets.SecretProvider, error) {
// 	factProvider, err := createFactProvider(conf.FactsProviderOptions)
// 	if err != nil {
// 	  return nil, nil, err
// 	}

// 	secretProvider, err := createSecretsProvider(conf.SecretsProviderOptions)
// 	if err != nil {
// 	  return nil, nil, err
// 	}

// 	return factProvider, secretProvider, nil
// }

// func createFactProvider(conf FactsProviderOptions) (facts.FactProvider, error) {
// 	switch conf.Provider {
// 	case "redis":
// 		opts, err := redis.ParseURL(conf.RedisURI)
// 		if err != nil {
// 		  return nil, err
// 		}

// 		return facts.NewRedisFactProvider(redis.NewClient(opts)), nil
// 	default:
// 		return nil, fmt.Errorf("no such fact provider %s", conf.Provider)
// 	}
// }

// func createSecretsProvider(conf SecretsProviderOptions) (secrets.SecretProvider, error) {
// 	switch conf.Provider {
// 	case "asm":
// 		cfgOpts := []func(*config.LoadOptions) error{}
// 		config, err := config.LoadDefaultConfig(context.Background(), cfgOpts...)
// 		if err != nil {
// 		  return nil, err
// 		}

// 		return secrets.NewASMSecretProvider(secretsmanager.NewFromConfig(config)), nil
// 	default:
// 		return nil, fmt.Errorf("no such secret provider %s", conf.Provider)
// 	}
// }