package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/templates"
)

func (r *UIRoutes) IdentityFactsDashboard(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "", templates.Dashboard(templates.FactsDashboard(ctx.Param("id"))))
}