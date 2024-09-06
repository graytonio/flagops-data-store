package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/facts"
)


func (r *APIRoutes) GetIdentityFacts(ctx *gin.Context) {
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

func (r *APIRoutes) GetIdentityFact(ctx *gin.Context) {
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

func (r *APIRoutes) SetIdentityFact(ctx *gin.Context) {
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

func (r *APIRoutes) DeleteIdentityFact(ctx *gin.Context) {
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