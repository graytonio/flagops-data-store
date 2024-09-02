package facts

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Facts map[string]string

var (
	ErrIdentityNotFound = errors.New("identity not found")
	ErrSecretNotFound = errors.New("fact not found")
)

// TODO Mock facts provider

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