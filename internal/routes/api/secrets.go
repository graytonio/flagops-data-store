package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-store/internal/secrets"
)


func (r *APIRoutes) GetIdentitySecrets(ctx *gin.Context) {
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

func (r *APIRoutes) GetIdentitySecret(ctx *gin.Context) {
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

func (r *APIRoutes) SetIdentitySecret(ctx *gin.Context) {
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

func (r *APIRoutes) DeleteIdentitySecret(ctx *gin.Context) {
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