package secrets

import (
	"errors"

	"github.com/gin-gonic/gin"
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
