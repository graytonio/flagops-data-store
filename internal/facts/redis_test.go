package facts_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-store/internal/facts"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func getRedisContainer(ctx context.Context) (testcontainers.Container, redis.UniversalClient, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		return nil, nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	return redisC, client, nil
}

func closeRedisContainer(t *testing.T, ctx context.Context, container testcontainers.Container) {
	if err := container.Terminate(ctx); err != nil {
		t.Fatalf("Could not stop redis: %s", err)
	}
}

func TestGetAllIdentities(t *testing.T) {
	ctx := context.Background()
	container, client, err := getRedisContainer(ctx)
	if err != nil {
		t.Fatalf("Could not start redis: %s", err)
	}
	defer closeRedisContainer(t, ctx, container)

	client.Set(ctx, "customer0:fact0", "foo", 0)
	client.Set(ctx, "customer0:fact2", "foo", 0)
	client.Set(ctx, "customer1:fact0", "foo", 0)
	client.Set(ctx, "customer1:fact3", "foo", 0)
	client.Set(ctx, "customer2:fact0", "foo", 0)

	provider := facts.NewRedisFactProvider(client)

	ids, err := provider.GetAllIdentities(&gin.Context{})
	if assert.NoError(t, err) {
		assert.ElementsMatch(t, []string{"customer0", "customer1", "customer2"}, ids)
	}
}

func TestSetIdentityFact(t *testing.T) {
	var tests = []struct {
		name          string
		id            string
		key           string
		value         string
		expectError   bool
		expectedKey   string
		expectedValue string
	}{
		{
			name:          "core function",
			id:            "customer0",
			key:           "fact0",
			value:         "foo",
			expectError:   false,
			expectedKey:   "customer0:fact0",
			expectedValue: "foo",
		},
		{
			name:        "missing id",
			id:          "",
			key:         "fact0",
			value:       "foo",
			expectError: true,
		},
		{
			name:        "missing key",
			id:          "customer0",
			key:         "",
			value:       "foo",
			expectError: true,
		},
		{
			name:        "missing value",
			id:          "customer0",
			key:         "fact0",
			value:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			container, client, err := getRedisContainer(ctx)
			if err != nil {
				t.Fatalf("Could not start redis: %s", err)
			}
			defer closeRedisContainer(t, ctx, container)

			provider := facts.NewRedisFactProvider(client)

			err = provider.SetIdentityFact(&gin.Context{}, tt.id, tt.key, tt.value)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, client.Get(ctx, tt.expectedKey).Val(), tt.expectedValue)
				}
			}
		})
	}

}

func TestGetIdentityFacts(t *testing.T) {

	var tests = []struct {
		name          string
		id            string
		setupDB       func(client redis.UniversalClient)
		expectError   bool
		expectedFacts facts.Facts
	}{
		{
			name:        "core functionality",
			id:          "customer0",
			expectError: false,
			expectedFacts: facts.Facts{
				"fact0": "foo",
				"fact1": "bar",
			},
			setupDB: func(client redis.UniversalClient) {
				client.Set(context.Background(), "customer0:fact0", "foo", 0)
				client.Set(context.Background(), "customer0:fact1", "bar", 0)
			},
		},
		{
			name:        "missing id",
			id:          "",
			expectError: true,
			setupDB:     func(client redis.UniversalClient) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			container, client, err := getRedisContainer(ctx)
			if err != nil {
				t.Fatalf("Could not start redis: %s", err)
			}
			defer closeRedisContainer(t, ctx, container)

			provider := facts.NewRedisFactProvider(client)

			tt.setupDB(client)

			facts, err := provider.GetIdentityFacts(&gin.Context{}, tt.id)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, facts, tt.expectedFacts)
				}
			}
		})
	}
}

func TestDeleteIdentityFact(t *testing.T) {
	var tests = []struct {
		name        string
		id          string
		key         string
		setupDB     func(client redis.UniversalClient)
		expectError bool
	}{
		{
			name:        "core functionality",
			id:          "customer0",
			key:         "fact0",
			expectError: false,
			setupDB: func(client redis.UniversalClient) {
				client.Set(context.Background(), "customer0:fact0", "foo", 0)
			},
		},
		{
			name:        "key to delete does not exist",
			id:          "customer0",
			key:         "fact0",
			expectError: false,
			setupDB: func(client redis.UniversalClient) {
			},
		},
		{
			name:        "missing id",
			id:          "",
			key:         "fact0",
			expectError: true,
			setupDB:     func(client redis.UniversalClient) {},
		},
		{
			name:        "missing key",
			id:          "customer0",
			key:         "",
			expectError: true,
			setupDB:     func(client redis.UniversalClient) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			container, client, err := getRedisContainer(ctx)
			if err != nil {
				t.Fatalf("Could not start redis: %s", err)
			}
			defer closeRedisContainer(t, ctx, container)

			provider := facts.NewRedisFactProvider(client)

			tt.setupDB(client)

			err = provider.DeleteIdentityFact(&gin.Context{}, tt.id, tt.key)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				_, err := client.Get(context.Background(), fmt.Sprintf("%s:%s", tt.id, tt.key)).Result()
				assert.ErrorIs(t, err, redis.Nil)
			}
		})
	}
}

func TestDeleteIdentity(t *testing.T) {
	var tests = []struct {
		name        string
		id          string
		setupDB     func(client redis.UniversalClient)
		expectError bool
	}{
		{
			name:        "core functionality",
			id:          "customer0",
			expectError: false,
			setupDB: func(client redis.UniversalClient) {
				client.Set(context.Background(), "customer0:fact0", "foo", 0)
			},
		},
		{
			name:        "key to delete does not exist",
			id:          "customer0",
			expectError: false,
			setupDB: func(client redis.UniversalClient) {
			},
		},
		{
			name:        "missing id",
			id:          "",
			expectError: true,
			setupDB:     func(client redis.UniversalClient) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			container, client, err := getRedisContainer(ctx)
			if err != nil {
				t.Fatalf("Could not start redis: %s", err)
			}
			defer closeRedisContainer(t, ctx, container)

			provider := facts.NewRedisFactProvider(client)

			tt.setupDB(client)

			err = provider.DeleteIdentity(&gin.Context{}, tt.id)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				keys, _, _ := client.Scan(context.Background(), 0, tt.id, 100).Result()
				assert.Empty(t, keys)
			}
		})
	}
}
