package secrets

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var _ SecretProvider = &ASMSecretProvider{}

const secretPrefix = "flagops-secret-"

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
	return fmt.Sprintf("%s%s", secretPrefix, id)
}

func (a *ASMSecretProvider) getLogEntry(ctx *gin.Context) *logrus.Entry {
	entry := logrus.WithFields(logrus.Fields{
		"caller_path": ctx.FullPath(),
		"provider":    "asm",
		"api":         "secrets",
	})

	if ctx.Param("id") != "" {
		entry = entry.WithField("id", ctx.Param("id"))
	}

	if ctx.Param("key") != "" {
		entry = entry.WithField("key", ctx.Param("key"))
	}

	return entry
}

// GetIdentitySecrets implements SecretProvider.
func (a *ASMSecretProvider) GetIdentitySecrets(ctx *gin.Context, id string) (Secrets, error) {
	log := a.getLogEntry(ctx)
	log.Debug("fetching identity secrets from provider")
	res, err := a.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(a.getIdentitySecretKey(id)),
	})
	if err != nil {
		var aerr *types.ResourceNotFoundException
		if !errors.As(err, &aerr) { // If some other error other than not found fail
			log.WithError(err).Error("could not fetch identity from provider")
			return nil, err
		}
		log.Error("identity not found in provider")
		return nil, ErrIdentityNotFound
	}

	results := Secrets{}
	err = json.NewDecoder(bytes.NewBufferString(*res.SecretString)).Decode(&results)
	if err != nil {
		log.WithError(err).Error("identity not found in provider")
		return nil, err
	}

	return results, nil
}

// SetIdentitySecret implements SecretProvider.
func (a *ASMSecretProvider) SetIdentitySecret(ctx *gin.Context, id string, key string, value string) error {
	log := a.getLogEntry(ctx)
	log.Debug("setting secret for identity")
	currentValues, err := a.GetIdentitySecrets(ctx, id)
	if err != nil {
		if !errors.Is(err, ErrIdentityNotFound) {
			return err
		}

		// TODO Handle case where secret is marked for deletion

		_, err := a.client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(a.getIdentitySecretKey(id)),
			SecretString: aws.String(fmt.Sprintf("{\"%s\":\"%s\"}\n", key, value)),
		})
		if err != nil {
			log.WithError(err).Error("could not create new identity")
			return err
		}

		return nil
	}
	currentValues[key] = value

	payload := bytes.NewBuffer([]byte{})
	err = json.NewEncoder(payload).Encode(currentValues)
	if err != nil {
		log.WithError(err).Error("malformed secret data")
		return err
	}

	_, err = a.client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(a.getIdentitySecretKey(id)),
		SecretString: aws.String(payload.String()),
	})
	if err != nil {
		log.WithError(err).Error("could not update identity secret")
		return err
	}

	return nil
}

// GetAllIdentities implements SecretProvider.
func (a *ASMSecretProvider) GetAllIdentities(ctx *gin.Context) ([]string, error) {
	log := a.getLogEntry(ctx)
	log.Debug("fetching all identities from provider")
	ids := []string{}
	token := ""

	for {
		res, err := a.client.ListSecrets(ctx, &secretsmanager.ListSecretsInput{
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
			log.WithError(err).Error("could not fetch page of results from provider")
			return nil, err
		}
		log.WithField("identities", len(res.SecretList)).Debug("fetched page of results from provider")
		
		for _, s := range res.SecretList {
			ids = append(ids, strings.TrimPrefix(*s.Name, secretPrefix))
		}

		if res.NextToken == nil {
			break
		}

		token = *res.NextToken
	}

	return ids, nil
}

// DeleteIdentity implements SecretProvider.
func (a *ASMSecretProvider) DeleteIdentity(ctx *gin.Context, id string) error {
	log := a.getLogEntry(ctx)
	log.Debug("deleting identity")
	_, err := a.client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
		SecretId:             aws.String(a.getIdentitySecretKey(id)),
		RecoveryWindowInDays: aws.Int64(7), // TODO Make configurable
	})
	if err != nil {
		log.WithError(err).Error("could not delete identity")
		return err
	}

	return nil
}

// DeleteIdentitySecret implements SecretProvider.
func (a *ASMSecretProvider) DeleteIdentitySecret(ctx *gin.Context, id string, key string) error {
	log := a.getLogEntry(ctx)
	log.Debug("deleting identity secret")
	currentValues, err := a.GetIdentitySecrets(ctx, id)
	if err != nil {
		return err
	}
	delete(currentValues, key)

	payload := bytes.NewBuffer([]byte{})
	err = json.NewEncoder(payload).Encode(currentValues)
	if err != nil {
		log.WithError(err).Error("malformed secret data")
		return err
	}

	_, err = a.client.UpdateSecret(ctx, &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(a.getIdentitySecretKey(id)),
		SecretString: aws.String(payload.String()),
	})
	if err != nil {
		log.WithError(err).Error("could not update identity secret")
		return err
	}

	return nil
}
