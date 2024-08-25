package secrets

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/aws-sdk-go/aws"
)

var _ SecretProvider = &ASMSecretProvider{}

// A secrets provider based on AWS Secrets Manager
type ASMSecretProvider struct {
	client *secretsmanager.Client
}

func NewASMSecretProvider(client *secretsmanager.Client) *ASMSecretProvider {
	return &ASMSecretProvider{
		client: client,
	}
}

func (a *ASMSecretProvider) getIdentitySecretKey(id string) string {
	return fmt.Sprintf("flagops-secret-%s", id)
}

// GetIdentitySecrets implements SecretProvider.
func (a *ASMSecretProvider) GetIdentitySecrets(id string) (Secrets, error) {
	res, err := a.client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(a.getIdentitySecretKey(id)),
	})
	if err != nil {
		var aerr *types.ResourceNotFoundException
		if !errors.As(err, &aerr) { // If some other error other than not found fail
			return nil, err
		}

		_, err := a.client.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
			Name: aws.String(a.getIdentitySecretKey(id)),
		})
		if err != nil {
		  return nil, err
		}

		return Secrets{}, nil
	}


	results := Secrets{}
	err = json.NewDecoder(bytes.NewBufferString(*res.SecretString)).Decode(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// SetIdentitySecret implements SecretProvider.
func (a *ASMSecretProvider) SetIdentitySecret(id string, key string, value string) error {
	currentValues, err := a.GetIdentitySecrets(id)
	if err != nil {
	  return err
	}
	currentValues[key] = value
	
	payload := bytes.NewBuffer([]byte{})
	err = json.NewEncoder(payload).Encode(currentValues)
	if err != nil {
	  return err
	}

	_, err = a.client.UpdateSecret(context.TODO(), &secretsmanager.UpdateSecretInput{
		SecretId: aws.String(a.getIdentitySecretKey(id)),
		SecretString: aws.String(payload.String()),
	})
	if err != nil {
	  return err
	}

	return nil
}

// GetAllIdentities implements SecretProvider.
func (a *ASMSecretProvider) GetAllIdentities() ([]string, error) {
	ids := []string{}
	token := ""

	// TODO Timeout context
	for {
		res, err := a.client.ListSecrets(context.Background(), &secretsmanager.ListSecretsInput{
			MaxResults: aws.Int32(150),
			Filters: []types.Filter{
				{
					Key: "name",
					Values: []string{
						"flagops-secret",
					},
				},
			},
			NextToken: aws.String(token),
		})
		if err != nil {
			return nil, err
		}

		for _, s := range res.SecretList {
			ids = append(ids, *s.Name)
		}

		token = *res.NextToken
		if token == "" {
			break
		}
	}

	return ids, nil
}

// DeleteIdentity implements SecretProvider.
func (a *ASMSecretProvider) DeleteIdentity(id string) error {
	_, err := a.client.DeleteSecret(context.TODO(), &secretsmanager.DeleteSecretInput{
		SecretId: aws.String(a.getIdentitySecretKey(id)),
	})
	if err != nil {
	  return err
	}

	return nil
}

// DeleteIdentitySecret implements SecretProvider.
func (a *ASMSecretProvider) DeleteIdentitySecret(id string, key string) error {
	currentValues, err := a.GetIdentitySecrets(id)
	if err != nil {
	  return err
	}
	delete(currentValues, key)

	payload := bytes.NewBuffer([]byte{})
	err = json.NewEncoder(payload).Encode(currentValues)
	if err != nil {
	  return err
	}

	_, err = a.client.UpdateSecret(context.TODO(), &secretsmanager.UpdateSecretInput{
		SecretId: aws.String(a.getIdentitySecretKey(id)),
		SecretString: aws.String(payload.String()),
	})
	if err != nil {
	  return err
	}

	return nil
}
