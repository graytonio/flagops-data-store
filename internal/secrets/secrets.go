package secrets

import "errors"

type Secrets map[string]string

var (
	ErrIdentityNotFound = errors.New("identity not found")
	ErrSecretNotFound = errors.New("secret not found")
)

type SecretProvider interface {
	// Returns a list of all available identities in the provider
	GetAllIdentities() ([]string, error)

	// Deletes all records belonging to identity
	DeleteIdentity(id string) error

	// Returns all facts belonging to the identity
	GetIdentitySecrets(id string) (Secrets, error)

	// Set the key for the given identity to the value
	SetIdentitySecret(id string, key string, value string) error

	// Deletes the key for the given identity
	DeleteIdentitySecret(id string, key string) error
}