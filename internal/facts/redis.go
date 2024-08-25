package facts

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
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

// GetAllIdentities implements FactProvider.
func (r *RedisFactProvider) GetAllIdentities() ([]string, error) {
	prefixSet := make(map[string]struct{})
	cursor := uint64(0)

	for {
		keys, nextCursor, err := r.client.Scan(context.TODO(), cursor, "*", 100).Result()
		if err != nil {
			return nil, err
		}

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
func (r *RedisFactProvider) GetIdentityFacts(id string) (Facts, error) {
	if id == "" {
		return nil, errors.New("id is blank")
	}
	
	result := map[string]string{}
	cursor := uint64(0)

	for {
		keys, nextCursor, err := r.client.Scan(context.TODO(), cursor, fmt.Sprintf("%s:*", id), 100).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			val, err := r.client.Get(context.TODO(), key).Result()
			if err != nil {
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

	return result, nil
}

// SetIdentityFact implements FactProvider.
func (r *RedisFactProvider) SetIdentityFact(id string, key string, value string) error {
	if id == "" {
		return errors.New("id is blank")
	}

	if key == "" {
		return errors.New("key is blank")
	}

	if value == "" {
		return errors.New("value is blank")
	}

	err := r.client.Set(context.TODO(), r.getIdentityFactPath(id, key), value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteIdentityFact implements FactProvider.
func (r *RedisFactProvider) DeleteIdentityFact(id string, key string) error {
	if id == "" {
		return errors.New("id is blank")
	}

	if key == "" {
		return errors.New("key is blank")
	}

	err := r.client.Del(context.TODO(), r.getIdentityFactPath(id, key)).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteIdentity implements FactProvider.
func (r *RedisFactProvider) DeleteIdentity(id string) error {
	if id == "" {
		return errors.New("id is blank")
	}
	
	iter := r.client.Scan(context.TODO(), 0, fmt.Sprintf("%s:*", id), 0).Iterator()

	for iter.Next(context.TODO()) {
		key := iter.Val()
		err := r.client.Del(context.TODO(), key).Err()
		if err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
