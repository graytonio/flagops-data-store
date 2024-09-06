package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)


type identieiesSupportedProviders struct{
	Facts bool `json:"facts"`
	Secrets bool `json:"secrets"`
}
type getAllIdentitiesResponse struct {
	Identities map[string]identieiesSupportedProviders `json:"identities"`
}

func (r *APIRoutes) GetAllIdentities(ctx *gin.Context) {
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

func (r *APIRoutes) DeleteIdentity(ctx *gin.Context) {
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