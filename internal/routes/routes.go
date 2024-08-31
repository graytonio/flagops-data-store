package routes

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/facts"
	"github.com/graytonio/flagops-data-storage/internal/secrets"
)

type Routes struct {
	FactProvider facts.FactProvider
	SecretProvider secrets.SecretProvider
}

type identieiesSupportedProviders struct{
	Facts bool `json:"facts"`
	Secrets bool `json:"secrets"`
}
type getAllIdentitiesResponse struct {
	Identities map[string]identieiesSupportedProviders `json:"identities"`
}

func (r *Routes) GetAllIdentities(ctx *gin.Context) {
	factsIds, err := r.FactProvider.GetAllIdentities(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	secretsIds, err := r.SecretProvider.GetAllIdentities(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ids := map[string]identieiesSupportedProviders{}
	for _, f := range factsIds {
		tmp := ids[f]
		tmp.Facts = true
		ids[f] = tmp
	}

	for _, f := range secretsIds {
		tmp := ids[f]
		tmp.Secrets = true
		ids[f] = tmp
	}

	ctx.JSON(http.StatusOK, getAllIdentitiesResponse{
		Identities: ids,
	})
}

func (r *Routes) DeleteIdentity(ctx *gin.Context) {
	identity := ctx.Param("id")
	if identity == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter must not be empty"))
		return
	}

	err := r.FactProvider.DeleteIdentity(ctx, identity)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = r.SecretProvider.DeleteIdentity(ctx, identity)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (r *Routes) GetIdentityFacts(ctx *gin.Context) {
	identity := ctx.Param("id")
	if identity == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter must not be empty"))
		return
	}

	identityFacts, err := r.FactProvider.GetIdentityFacts(ctx, identity)
	if err != nil {
		if errors.Is(err, facts.ErrIdentityNotFound) {
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, identityFacts)
}

func (r *Routes) GetIdentityFact(ctx *gin.Context) {
	identity := ctx.Param("id")
	if identity == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter must not be empty"))
		return
	}

	key := ctx.Param("fact")
	if key == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("fact parameter must not be empty"))
		return
	}

	facts, err := r.FactProvider.GetIdentityFacts(ctx, identity)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	v, ok := facts[key]
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{
		key: v,
	})
}

type setIdentityFactRequest struct {
	Value string `json:"value"`
}

func (r *Routes) SetIdentityFact(ctx *gin.Context) {
	identity := ctx.Param("id")
	if identity == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter must not be empty"))
		return
	}

	key := ctx.Param("fact")
	if key == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("fact parameter must not be empty"))
		return
	}

	var body setIdentityFactRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := r.FactProvider.SetIdentityFact(ctx, identity, key, body.Value)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (r *Routes) DeleteIdentityFact(ctx *gin.Context) {
	identity := ctx.Param("id")
	if identity == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter must not be empty"))
		return
	}

	key := ctx.Param("fact")
	if key == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("fact parameter must not be empty"))
		return
	}

	err := r.FactProvider.DeleteIdentityFact(ctx, identity, key)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (r *Routes) GetIdentitySecrets(ctx *gin.Context) {
	identity := ctx.Param("id")
	if identity == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter must not be empty"))
		return
	}

	facts, err := r.SecretProvider.GetIdentitySecrets(ctx, identity)
	if err != nil {
		if errors.Is(err, secrets.ErrIdentityNotFound) {
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}
		
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, facts)
}

func (r *Routes) GetIdentitySecret(ctx *gin.Context) {
	identity := ctx.Param("id")
	if identity == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter must not be empty"))
		return
	}

	key := ctx.Param("secret")
	if key == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("secret parameter must not be empty"))
		return
	}

	identitySecrets, err := r.SecretProvider.GetIdentitySecrets(ctx, identity)
	if err != nil {
		if errors.Is(err, secrets.ErrIdentityNotFound) {
			ctx.AbortWithError(http.StatusNotFound, err)
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	v, ok := identitySecrets[key]
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{
		key: v,
	})
}

type setIdentitySecretRequest struct {
	Value string `json:"value"`
}

func (r *Routes) SetIdentitySecret(ctx *gin.Context) {
	identity := ctx.Param("id")
	if identity == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter must not be empty"))
		return
	}

	key := ctx.Param("secret")
	if key == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("secret parameter must not be empty"))
		return
	}

	var body setIdentitySecretRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := r.SecretProvider.SetIdentitySecret(ctx, identity, key, body.Value)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (r *Routes) DeleteIdentitySecret(ctx *gin.Context) {
	identity := ctx.Param("id")
	if identity == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter must not be empty"))
		return
	}

	key := ctx.Param("secret")
	if key == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("secret parameter must not be empty"))
		return
	}

	err := r.SecretProvider.DeleteIdentitySecret(ctx, identity, key)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}