package facts

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var _ FactProvider = &RedisFactProvider{}

type RedisFactProvider struct {
	client redis.UniversalClient
}

func NewRedisFactProvider(client redis.UniversalClient) *RedisFactProvider {
	return &RedisFactProvider{
		client: client,
	}
}

func (r *RedisFactProvider) getIdentityFactPath(id string, key string) string {
	return fmt.Sprintf("%s:%s", id, key)
}

func (r *RedisFactProvider) getLogEntry(ctx *gin.Context) *logrus.Entry {
	entry := logrus.WithFields(logrus.Fields{
		"caller_path": ctx.FullPath(),
		"provider":    "redis",
		"api":         "facts",
	})

	if ctx.Param("id") != "" {
		entry = entry.WithField("id", ctx.Param("id"))
	}

	if ctx.Param("key") != "" {
		entry = entry.WithField("key", ctx.Param("key"))
	}

	return entry
}

// GetAllIdentities implements FactProvider.
func (r *RedisFactProvider) GetAllIdentities(ctx *gin.Context) ([]string, error) {
	log := r.getLogEntry(ctx)
	log.Debug("fetching all identities from provider")
	prefixSet := make(map[string]struct{})
	cursor := uint64(0)

	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, "*", 100).Result()
		if err != nil {
			log.WithError(err).Error("could not fetch scan page from provider")
			return nil, err
		}
		log.WithField("identities", len(keys)).Debug("fetched page of identities")

		for _, key := range keys {
			parts := strings.SplitN(key, ":", 2)
			if len(parts) > 1 {
				prefixSet[parts[0]] = struct{}{}
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	prefixes := make([]string, 0, len(prefixSet))
	for prefix := range prefixSet {
		prefixes = append(prefixes, prefix)
	}

	return prefixes, nil
}

// GetIdentityFacts implements FactProvider.
func (r *RedisFactProvider) GetIdentityFacts(ctx *gin.Context, id string) (Facts, error) {
	log := r.getLogEntry(ctx)
	if id == "" {
		log.Debug("called with no identity")
		return nil, errors.New("id is blank")
	}

	result := map[string]string{}
	cursor := uint64(0)

	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, fmt.Sprintf("%s:*", id), 100).Result()
		if err != nil {
			log.WithError(err).Error("could not fetch scan page from provider")
			return nil, err
		}
		log.WithField("keys", len(keys)).Debug("fetched page of keys")

		for _, key := range keys {
			val, err := r.client.Get(ctx, key).Result()
			if err != nil {
				log.WithError(err).Error("could not fetch key from provider")
				return nil, err
			}
			parts := strings.SplitN(key, ":", 2)
			result[parts[1]] = val
		}

		// Move the cursor to the next batch
		cursor = nextCursor

		// If cursor is 0, the iteration is complete
		if cursor == 0 {
			break
		}
	}

	if len(result) == 0 {
		return nil, ErrIdentityNotFound
	}

	return result, nil
}

// SetIdentityFact implements FactProvider.
func (r *RedisFactProvider) SetIdentityFact(ctx *gin.Context, id string, key string, value string) error {
	log := r.getLogEntry(ctx)
	log.Debug("setting fact for identity")
	if id == "" {
		log.Debug("called with no identity")
		return errors.New("id is blank")
	}

	if key == "" {
		log.Debug("called with no key")
		return errors.New("key is blank")
	}

	if value == "" {
		log.Debug("called with no value")
		return errors.New("value is blank")
	}

	log.Debug("setting identity fact")
	err := r.client.Set(ctx, r.getIdentityFactPath(id, key), value, 0).Err()
	if err != nil {
		log.WithError(err).Error("could not set key in provider")
		return err
	}

	return nil
}

// DeleteIdentityFact implements FactProvider.
func (r *RedisFactProvider) DeleteIdentityFact(ctx *gin.Context, id string, key string) error {
	log := r.getLogEntry(ctx)
	if id == "" {
		log.Debug("called with no identity")
		return errors.New("id is blank")
	}

	if key == "" {
		log.Debug("called with no key")
		return errors.New("key is blank")
	}

	log.Debug("deleting identity fact")
	err := r.client.Del(ctx, r.getIdentityFactPath(id, key)).Err()
	if err != nil {
		log.WithError(err).Error("could not delete key in provider")
		return err
	}

	return nil
}

// DeleteIdentity implements FactProvider.
func (r *RedisFactProvider) DeleteIdentity(ctx *gin.Context, id string) error {
	log := r.getLogEntry(ctx)
	if id == "" {
		log.Debug("called with no identity")
		return errors.New("id is blank")
	}

	cursor := uint64(0)
	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, fmt.Sprintf("%s:*", id), 100).Result()
		if err != nil {
			log.WithError(err).Error("could not fetch scan page from provider")
			return err
		}
		log.WithField("keys", len(keys)).Debug("fetched page of keys")

		for _, key := range keys {
			_, err := r.client.Del(ctx, key).Result()
			if err != nil {
				log.WithError(err).Error("could not delete key in provider")
				return err
			}
		}

		// Move the cursor to the next batch
		cursor = nextCursor

		// If cursor is 0, the iteration is complete
		if cursor == 0 {
			break
		}
	}

	return nil
}
