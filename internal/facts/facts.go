package facts

type Facts map[string]string

type FactProvider interface {
	// Returns a list of all available identities in the provider
	GetAllIdentities() ([]string, error)

	// Deletes all records belonging to identity
	DeleteIdentity(id string) error

	// Returns all facts belonging to the identity
	GetIdentityFacts(id string) (Facts, error)

	// Set the key for the given identity to the value
	SetIdentityFact(id string, key string, value string) error

	// Deletes the key for the given identity
	DeleteIdentityFact(id string, key string) error
}