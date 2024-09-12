package ui

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-store/internal/facts"
	"github.com/graytonio/flagops-data-store/internal/utils"
	"github.com/graytonio/flagops-data-store/templates/components"
	"github.com/graytonio/flagops-data-store/templates/pages"
)

func SendHTMXError(ctx *gin.Context, code int, message string) {
	ctx.Header("HX-Retarget", "#error-title")
	ctx.HTML(code, "", components.ErrorTitle(message))
}

func (r *UIRoutes) IdentityFactsTable(ctx *gin.Context) {
	identityFacts, err := r.FactProvider.GetIdentityFacts(ctx, ctx.Param("id"))
	if err != nil {
		if errors.Is(err, facts.ErrIdentityNotFound) {
			SendHTMXError(ctx, http.StatusNotFound, fmt.Sprintf("%s not found", ctx.Param("id")))
			return
		}
		SendHTMXError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "", pages.IdentityFacts(pages.IdentityFactsViewData{
		Identity: ctx.Param("id"),
		Facts:    identityFacts,
	}))
}

func (r *UIRoutes) EditIdentityFactRowForm(ctx *gin.Context) {
	id := ctx.Param("id")
	fact := ctx.Param("fact")
	
	facts, err := r.FactProvider.GetIdentityFacts(ctx, id)
	if err != nil {
		SendHTMXError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "", pages.EditRow(id, fact, facts[fact]))
}

type identityFactEdit struct {
	NewValue string `form:"value"`
}

func (r *UIRoutes) EditIdentityFactRow(ctx *gin.Context) {
	id := ctx.Param("id")
	fact := ctx.Param("fact")

	var data identityFactEdit

	if err := ctx.Bind(&data); err != nil {
		SendHTMXError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if err := r.FactProvider.SetIdentityFact(ctx, id, fact, data.NewValue); err != nil {
		SendHTMXError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "", pages.FactRow(id, fact, data.NewValue))
}

type identitySearchRequest struct {
	SearchData string `form:"search"`
}

func (r *UIRoutes) IdentitySearch(ctx *gin.Context) {
	var searchData identitySearchRequest
	err := ctx.Bind(&searchData)
	if err != nil {
		SendHTMXError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	factsIds, err := r.FactProvider.GetAllIdentities(ctx)
	if err != nil {
		SendHTMXError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	secretsIds, err := r.SecretProvider.GetAllIdentities(ctx)
	if err != nil {
		SendHTMXError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	identities := utils.RemoveDuplicate(append(factsIds, secretsIds...))

	searchResults := slices.DeleteFunc(identities, func(id string) bool {
		return !strings.Contains(id, searchData.SearchData) // TODO Better searching
	})

	ctx.HTML(http.StatusOK, "", pages.IdentitiesSearchResults(searchResults))
}
