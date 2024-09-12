package secrets

import (
	"context"
	"errors"
	"fmt"

	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-store/internal/config"
)

type Secrets map[string]string

var (
	ErrIdentityNotFound = errors.New("identity not found")
	ErrSecretNotFound   = errors.New("secret not found")
)

type SecretProvider interface {
	// Returns a list of all available identities in the provider
	GetAllIdentities(ctx *gin.Context) ([]string, error)

	// Deletes all records belonging to identity
	DeleteIdentity(ctx *gin.Context, id string) error

	// Returns all Secrets belonging to the identity
	GetIdentitySecrets(ctx *gin.Context, id string) (Secrets, error)

	// Set the key for the given identity to the value
	SetIdentitySecret(ctx *gin.Context, id string, key string, value string) error

	// Deletes the key for the given identity
	DeleteIdentitySecret(ctx *gin.Context, id string, key string) error
}

func GetSecretsProvider(conf config.SecretsProviderOptions) (SecretProvider, error) {
	switch conf.Provider {
	case "asm":
		config, err := awsconf.LoadDefaultConfig(context.Background())
		if err != nil {
			return nil, err
		}

		return NewASMSecretProvider(secretsmanager.NewFromConfig(config), conf), nil
	default:
		return nil, fmt.Errorf("no such secret provider %s", conf.Provider)
	}
}

var _ SecretProvider = &MockSecretsProvider{}

type MockSecretsProvider struct {
	SecretsDB map[string]map[string]string // Holds our "Secrets lookup table"
}

// DeleteIdentity implements FactProvider.
func (m *MockSecretsProvider) DeleteIdentity(ctx *gin.Context, id string) error {
	if _, ok := m.SecretsDB[id]; !ok {
		return nil
	}

	delete(m.SecretsDB, id)
	return nil
}

// DeleteIdentityFact implements FactProvider.
func (m *MockSecretsProvider) DeleteIdentitySecret(ctx *gin.Context, id string, key string) error {
	if _, ok := m.SecretsDB[id]; !ok {
		return nil
	}

	if _, ok := m.SecretsDB[id][key]; !ok {
		return nil
	}

	delete(m.SecretsDB[id], key)
	return nil
}

// GetAllIdentities implements FactProvider.
func (m *MockSecretsProvider) GetAllIdentities(ctx *gin.Context) ([]string, error) {
	ids := []string{}
	for k := range m.SecretsDB {
		ids = append(ids, k)
	}
	return ids, nil
}

// GetIdentitySecrets implements FactProvider.
func (m *MockSecretsProvider) GetIdentitySecrets(ctx *gin.Context, id string) (Secrets, error) {
	identitySecrets, ok := m.SecretsDB[id]
	if !ok {
		return nil, ErrIdentityNotFound
	}

	return identitySecrets, nil
}

// SetIdentityFact implements FactProvider.
func (m *MockSecretsProvider) SetIdentitySecret(ctx *gin.Context, id string, key string, value string) error {
	identitySecrets, ok := m.SecretsDB[id]
	if !ok {
		return ErrIdentityNotFound
	}

	identitySecrets[key] = value
	return nil
}
