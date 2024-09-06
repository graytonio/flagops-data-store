package ui

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/facts"
	"github.com/graytonio/flagops-data-storage/templates"
)

func (r *UIRoutes) IdentityFactsTable(ctx *gin.Context) {
	identityFacts, err := r.FactProvider.GetIdentityFacts(ctx, ctx.Param("id"))
	if err != nil {
		if errors.Is(err, facts.ErrIdentityNotFound) {
			ctx.Header("HX-Retarget", "#identity-title")
			ctx.HTML(http.StatusNotFound, "", templates.IdentityNotFoundError(ctx.Param("id"))) // TODO Change to html
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err) // TODO Change to html
		return
	}

	ctx.HTML(http.StatusOK, "", templates.FactsTable(templates.FactsDashboardViewData{
		Identity: ctx.Param("id"),
		Facts: identityFacts,
	}))
}