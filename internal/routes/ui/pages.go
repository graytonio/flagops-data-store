package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/utils"
	"github.com/graytonio/flagops-data-storage/templates"
	"github.com/graytonio/flagops-data-storage/templates/layout"
	"github.com/graytonio/flagops-data-storage/templates/pages"
)

func (r *UIRoutes) HomeDashboard(ctx *gin.Context) {
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

	identities := utils.RemoveDuplicate(append(factsIds, secretsIds...))
	
	ctx.HTML(http.StatusOK, "", layout.DashboardLayout(
		pages.DashboardContent(pages.DashboardViewData{
			Identities: identities,
		}),
	))
}

func (r *UIRoutes) IdentityFactsDashboard(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "", layout.DashboardLayout(templates.FactsDashboard(ctx.Param("id"))))
}