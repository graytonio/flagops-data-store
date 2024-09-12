package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/templates/layout"
	"github.com/graytonio/flagops-data-storage/templates/pages"
)

func (r *UIRoutes) HomeDashboard(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "", layout.Layout(layout.DashboardLayout(
		pages.IdentitiesPage(),
	)))
}

func (r *UIRoutes) IdentityFactsDashboard(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "", layout.Layout(layout.DashboardLayout(
		pages.IdentityDetailsPage(ctx.Param("id")),
	)))
}
