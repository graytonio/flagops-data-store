package facts

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/config"
	"github.com/redis/go-redis/v9"
)

type Facts map[string]string

var (
	ErrIdentityNotFound = errors.New("identity not found")
	ErrSecretNotFound   = errors.New("fact not found")
)

type FactProvider interface {
	// Returns a list of all available identities in the provider
	GetAllIdentities(ctx *gin.Context) ([]string, error)

	// Deletes all records belonging to identity
	DeleteIdentity(ctx *gin.Context, id string) error

	// Returns all facts belonging to the identity
	GetIdentityFacts(ctx *gin.Context, id string) (Facts, error)

	// Set the key for the given identity to the value
	SetIdentityFact(ctx *gin.Context, id string, key string, value string) error

	// Deletes the key for the given identity
	DeleteIdentityFact(ctx *gin.Context, id string, key string) error
}

func GetFactProvider(config config.FactsProviderOptions) (FactProvider, error) {
	switch config.Provider {
	case "redis":
		opts, err := redis.ParseURL(config.RedisURI)
		if err != nil {
		  return nil, err
		}

		return NewRedisFactProvider(redis.NewClient(opts)), nil
	default:
		return nil, fmt.Errorf("no such fact provider %s", config.Provider)
	}
}

var _ FactProvider = &MockFactsProvider{}

type MockFactsProvider struct {
	FactsDB map[string]map[string]string // Holds our "facts lookup table"
}

// DeleteIdentity implements FactProvider.
func (m *MockFactsProvider) DeleteIdentity(ctx *gin.Context, id string) error {
	if _, ok := m.FactsDB[id]; !ok {
		return nil
	}
	
	delete(m.FactsDB, id)
	return nil
}

// DeleteIdentityFact implements FactProvider.
func (m *MockFactsProvider) DeleteIdentityFact(ctx *gin.Context, id string, key string) error {
	if _, ok := m.FactsDB[id]; !ok {
		return nil
	}

	if _, ok := m.FactsDB[id][key]; !ok {
		return nil
	}
	
	delete(m.FactsDB[id], key)
	return nil
}

// GetAllIdentities implements FactProvider.
func (m *MockFactsProvider) GetAllIdentities(ctx *gin.Context) ([]string, error) {
	ids := []string{}
	for k := range m.FactsDB {
		ids = append(ids, k)
	}
	return ids, nil
}

// GetIdentityFacts implements FactProvider.
func (m *MockFactsProvider) GetIdentityFacts(ctx *gin.Context, id string) (Facts, error) {
	identityFacts, ok := m.FactsDB[id]
	if !ok {
		return nil, ErrIdentityNotFound
	}

	return identityFacts, nil
}

// SetIdentityFact implements FactProvider.
func (m *MockFactsProvider) SetIdentityFact(ctx *gin.Context, id string, key string, value string) error {
	identityFacts, ok := m.FactsDB[id]
	if !ok {
		return ErrIdentityNotFound
	}

	identityFacts[key] = value
	return nil
}
