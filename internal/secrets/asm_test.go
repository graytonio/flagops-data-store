package secrets_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/secrets"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

func getLocalStackContainer(ctx context.Context, region string) (testcontainers.Container, *secretsmanager.Client, error) {
	stackC, err := localstack.Run(ctx, "localstack/localstack:3.7.0")
	if err != nil {
		return nil, nil, err
	}

	provider, err := testcontainers.NewDockerProvider()
    if err != nil {
        return nil, nil, err
    }
    defer provider.Close()

    host, err := provider.DaemonHost(ctx)
    if err != nil {
        return nil, nil, err
    }

	mappedPort, err := stackC.MappedPort(ctx, nat.Port("4566/tcp"))
    if err != nil {
        return nil, nil, err
    }

	endpoint := fmt.Sprintf("http://%s:%d", host, mappedPort.Int())

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("accessKey", "secretKey", "token")),
	)
	if err != nil {
		return nil, nil, err
	}

	client := secretsmanager.NewFromConfig(awsCfg, func(o *secretsmanager.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return stackC, client, nil
}

func closeLocalStackContainer(t *testing.T, ctx context.Context, container testcontainers.Container) {
	if err := container.Terminate(ctx); err != nil {
		t.Fatalf("Could not stop local stack: %s", err)
	}
}

func TestGetIdentitySecrets(t *testing.T) {
	ctx := context.Background()
	container, client, err := getLocalStackContainer(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("Could not start local stack: %s", err)
	}
	defer closeLocalStackContainer(t, ctx, container)

	client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String(fmt.Sprintf("flagops-secret-%s", "test-identity")),
		SecretString: aws.String(`{"foo": "bar", "boo": "baz"}`),
	})

	provider := secrets.NewASMSecretProvider(client)

	secretsOutput, err := provider.GetIdentitySecrets(&gin.Context{}, "test-identity")
	if assert.NoError(t, err) {
		assert.Equal(t, secrets.Secrets{"foo": "bar", "boo": "baz"}, secretsOutput)
	}
}

func TestSetIdentitySecretNoExistingSecret(t *testing.T) {
	ctx := context.Background()
	container, client, err := getLocalStackContainer(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("Could not start local stack: %s", err)
	}
	defer closeLocalStackContainer(t, ctx, container)

	provider := secrets.NewASMSecretProvider(client)

	err = provider.SetIdentitySecret(&gin.Context{}, "test-identity", "foo", "bar")
	if assert.NoError(t, err) {
		res, _ := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(fmt.Sprintf("flagops-secret-%s", "test-identity")),
		})

		assert.Equal(t, "{\"foo\":\"bar\"}\n", *res.SecretString)
	}
}

func TestSetIdentitySecretExistingSecret(t *testing.T) {
	ctx := context.Background()
	container, client, err := getLocalStackContainer(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("Could not start local stack: %s", err)
	}
	defer closeLocalStackContainer(t, ctx, container)

	client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String(fmt.Sprintf("flagops-secret-%s", "test-identity")),
		SecretString: aws.String(`{"foo":"baz"}`),
	})

	provider := secrets.NewASMSecretProvider(client)

	err = provider.SetIdentitySecret(&gin.Context{}, "test-identity", "foo", "bar")
	if assert.NoError(t, err) {
		res, _ := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(fmt.Sprintf("flagops-secret-%s", "test-identity")),
		})

		assert.Equal(t, "{\"foo\":\"bar\"}\n", *res.SecretString)
	}
}

func TestGetAllIdenties(t *testing.T) {
	ctx := context.Background()
	container, client, err := getLocalStackContainer(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("Could not start local stack: %s", err)
	}
	defer closeLocalStackContainer(t, ctx, container)

	client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String(fmt.Sprintf("flagops-secret-%s", "test-identity-0")),
		SecretString: aws.String(`{"foo":"baz"}`),
	})

	client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String(fmt.Sprintf("flagops-secret-%s", "test-identity-1")),
		SecretString: aws.String(`{"foo":"baz"}`),
	})

	client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String(fmt.Sprintf("flagops-secret-%s", "test-identity-2")),
		SecretString: aws.String(`{"foo":"baz"}`),
	})

	provider := secrets.NewASMSecretProvider(client)

	ids, err := provider.GetAllIdentities(&gin.Context{})
	if assert.NoError(t, err) {
		assert.ElementsMatch(t, []string{"test-identity-0", "test-identity-1", "test-identity-2"}, ids)
	}
}

func TestDeleteIdentitySecret(t *testing.T) {
	ctx := context.Background()
	container, client, err := getLocalStackContainer(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("Could not start local stack: %s", err)
	}
	defer closeLocalStackContainer(t, ctx, container)

	client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String(fmt.Sprintf("flagops-secret-%s", "test-identity")),
		SecretString: aws.String(`{"foo":"baz","boo":"bar"}`),
	})

	provider := secrets.NewASMSecretProvider(client)

	err = provider.DeleteIdentitySecret(&gin.Context{}, "test-identity", "boo")
	if assert.NoError(t, err) {
		res, _ := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(fmt.Sprintf("flagops-secret-%s", "test-identity")),
		})

		assert.Equal(t, "{\"foo\":\"baz\"}\n", *res.SecretString)
	}
}

func TestDeleteIdentity(t *testing.T) {
	ctx := context.Background()
	container, client, err := getLocalStackContainer(ctx, "us-east-1")
	if err != nil {
		t.Fatalf("Could not start local stack: %s", err)
	}
	defer closeLocalStackContainer(t, ctx, container)

	client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String(fmt.Sprintf("flagops-secret-%s", "test-identity")),
		SecretString: aws.String(`{"foo":"baz","boo":"bar"}`),
	})

	provider := secrets.NewASMSecretProvider(client)

	err = provider.DeleteIdentity(&gin.Context{}, "test-identity")
	if assert.NoError(t, err) {
		_, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(fmt.Sprintf("flagops-secret-%s", "test-identity")),
		})
		
		var aerr *types.InvalidRequestException
		assert.ErrorAs(t, err, &aerr)
	}
}