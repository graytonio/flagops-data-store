package flagopsdatastore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/open-feature/go-sdk/openfeature"
)

type ProviderOption func(*Provider)

var _ openfeature.FeatureProvider = &Provider{}

type Provider struct {
	httpClient      *http.Client
	baseURL         *url.URL
	featureProvider openfeature.FeatureProvider
}

func NewProvider(rawBaseURL string, featureProvider openfeature.FeatureProvider) (*Provider, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
	  return nil, err
	}
	
	provider := &Provider{
		httpClient:      &http.Client{},
		baseURL: baseURL,
		featureProvider: featureProvider,
	}

	return provider, nil
}

// Metadata implements openfeature.FeatureProvider.
func (p *Provider) Metadata() openfeature.Metadata {
	return openfeature.Metadata{
		Name: "flagops-data-store",
	}
}

// Hooks implements openfeature.FeatureProvider.
func (p *Provider) Hooks() []openfeature.Hook {
	return []openfeature.Hook{}
}

func (p *Provider) getIdentityContext(id string) (map[string]string, error) {
	reqURL := p.baseURL.JoinPath("/fact", id)

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
	  return nil, err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
	  return nil, err
	}
	defer resp.Body.Close()

	facts := map[string]string{}
	err = json.NewDecoder(resp.Body).Decode(&facts)
	if err != nil {
	  return nil, err
	}

	return facts, nil
}

func injectIdentityContext(identityCtx map[string]string, evalCtx openfeature.FlattenedContext) openfeature.FlattenedContext {
	for k, v := range identityCtx {
		evalCtx[k] = v
	}

	return evalCtx
}

// BooleanEvaluation implements openfeature.FeatureProvider.
func (p *Provider) BooleanEvaluation(ctx context.Context, flag string, defaultValue bool, evalCtx openfeature.FlattenedContext) openfeature.BoolResolutionDetail {
	targetKey, ok := evalCtx[openfeature.TargetingKey]
	if !ok {
		return p.featureProvider.BooleanEvaluation(ctx, flag, defaultValue, evalCtx)
	}

	context, err := p.getIdentityContext(targetKey.(string))
	if err != nil {
	  return openfeature.BoolResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			ResolutionError: openfeature.NewGeneralResolutionError("could not fetch identity context from data store"),
		},
	  }
	}

	evalCtx = injectIdentityContext(context, evalCtx)
	return p.featureProvider.BooleanEvaluation(ctx, flag, defaultValue, evalCtx)
}

// FloatEvaluation implements openfeature.FeatureProvider.
func (p *Provider) FloatEvaluation(ctx context.Context, flag string, defaultValue float64, evalCtx openfeature.FlattenedContext) openfeature.FloatResolutionDetail {
	targetKey, ok := evalCtx[openfeature.TargetingKey]
	if !ok {
		return p.featureProvider.FloatEvaluation(ctx, flag, defaultValue, evalCtx)
	}

	context, err := p.getIdentityContext(targetKey.(string))
	if err != nil {
	  return openfeature.FloatResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			ResolutionError: openfeature.NewGeneralResolutionError("could not fetch identity context from data store"),
		},
	  }
	}

	evalCtx = injectIdentityContext(context, evalCtx)
	return p.featureProvider.FloatEvaluation(ctx, flag, defaultValue, evalCtx)
}

// IntEvaluation implements openfeature.FeatureProvider.
func (p *Provider) IntEvaluation(ctx context.Context, flag string, defaultValue int64, evalCtx openfeature.FlattenedContext) openfeature.IntResolutionDetail {
	targetKey, ok := evalCtx[openfeature.TargetingKey]
	if !ok {
		return p.featureProvider.IntEvaluation(ctx, flag, defaultValue, evalCtx)
	}

	context, err := p.getIdentityContext(targetKey.(string))
	if err != nil {
	  return openfeature.IntResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			ResolutionError: openfeature.NewGeneralResolutionError("could not fetch identity context from data store"),
		},
	  }
	}

	evalCtx = injectIdentityContext(context, evalCtx)
	return p.featureProvider.IntEvaluation(ctx, flag, defaultValue, evalCtx)
}

// ObjectEvaluation implements openfeature.FeatureProvider.
func (p *Provider) ObjectEvaluation(ctx context.Context, flag string, defaultValue interface{}, evalCtx openfeature.FlattenedContext) openfeature.InterfaceResolutionDetail {
	targetKey, ok := evalCtx[openfeature.TargetingKey]
	if !ok {
		return p.featureProvider.ObjectEvaluation(ctx, flag, defaultValue, evalCtx)
	}

	context, err := p.getIdentityContext(targetKey.(string))
	if err != nil {
	  return openfeature.InterfaceResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			ResolutionError: openfeature.NewGeneralResolutionError("could not fetch identity context from data store"),
		},
	  }
	}

	evalCtx = injectIdentityContext(context, evalCtx)
	return p.featureProvider.ObjectEvaluation(ctx, flag, defaultValue, evalCtx)
}

// StringEvaluation implements openfeature.FeatureProvider.
func (p *Provider) StringEvaluation(ctx context.Context, flag string, defaultValue string, evalCtx openfeature.FlattenedContext) openfeature.StringResolutionDetail {
	targetKey, ok := evalCtx[openfeature.TargetingKey]
	if !ok {
		return p.featureProvider.StringEvaluation(ctx, flag, defaultValue, evalCtx)
	}

	context, err := p.getIdentityContext(targetKey.(string))
	if err != nil {
	  return openfeature.StringResolutionDetail{
		Value: defaultValue,
		ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
			ResolutionError: openfeature.NewGeneralResolutionError("could not fetch identity context from data store"),
		},
	  }
	}

	evalCtx = injectIdentityContext(context, evalCtx)
	return p.featureProvider.StringEvaluation(ctx, flag, defaultValue, evalCtx)
}

func (p *Provider) GetIdentityFacts(ctx context.Context, identity string) (map[string]string, error) {
	reqURL := p.baseURL.JoinPath("/fact", identity)

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
	  return nil, err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
	  return nil, err
	}
	defer resp.Body.Close()

	facts := map[string]string{}
	err = json.NewDecoder(resp.Body).Decode(&facts)
	if err != nil {
	  return nil, err
	}

	return facts, nil
}

func (p *Provider) SetIdentityFact(ctx context.Context, identity string, key string, value string) error {
	reqUrl := p.baseURL.JoinPath("/fact", identity, key)

	body := map[string]string{
		"value": value,
	}

	bodyBytes := bytes.NewBuffer(nil)
	err := json.NewEncoder(bodyBytes).Encode(body)
	if err != nil {
	  return err
	}

	req, err := http.NewRequest(http.MethodPut, reqUrl.String(), bodyBytes)
	if err != nil {
	  return err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
	  return err
	}

	if resp.StatusCode == 200 {
		return nil
	}

	return errors.New(resp.Status)
}

func (p *Provider) GetIdentitySecrets(ctx context.Context, identity string) (map[string]string, error) {
	reqURL := p.baseURL.JoinPath("/secret", identity)

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
	  return nil, err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
	  return nil, err
	}
	defer resp.Body.Close()

	secrets := map[string]string{}
	err = json.NewDecoder(resp.Body).Decode(&secrets)
	if err != nil {
	  return nil, err
	}

	return secrets, nil
}

func (p *Provider) SetIdentitySecret(ctx context.Context, identity string, key string, value string) error {
	reqUrl := p.baseURL.JoinPath("/secret", identity, key)

	body := map[string]string{
		"value": value,
	}

	bodyBytes := bytes.NewBuffer(nil)
	err := json.NewEncoder(bodyBytes).Encode(body)
	if err != nil {
	  return err
	}

	req, err := http.NewRequest(http.MethodPut, reqUrl.String(), bodyBytes)
	if err != nil {
	  return err
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
	  return err
	}

	if resp.StatusCode == 200 {
		return nil
	}

	return errors.New(resp.Status)
}