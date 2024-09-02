package secrets

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Secrets map[string]string

var (
	ErrIdentityNotFound = errors.New("identity not found")
	ErrSecretNotFound = errors.New("secret not found")
)

// Mock sercrets provider

type SecretProvider interface {
	// Returns a list of all available identities in the provider
	GetAllIdentities(ctx *gin.Context) ([]string, error)

	// Deletes all records belonging to identity
	DeleteIdentity(ctx *gin.Context, id string) error

	// Returns all facts belonging to the identity
	GetIdentitySecrets(ctx *gin.Context, id string) (Secrets, error)

	// Set the key for the given identity to the value
	SetIdentitySecret(ctx *gin.Context, id string, key string, value string) error

	// Deletes the key for the given identity
	DeleteIdentitySecret(ctx *gin.Context, id string, key string) error
}