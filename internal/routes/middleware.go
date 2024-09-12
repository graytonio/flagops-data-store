package routes

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-store/internal/db"
	"github.com/sirupsen/logrus"
)

func ErrorLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		for _, ginErr := range ctx.Errors {
			logrus.WithError(ginErr).Error("error handling http request")
		}
	}
}

func (r *Routes) RequiresAuth(permissions ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !r.Config.UserDatabaseOptions.RequireAuth {
			ctx.Next()
			return
		}
		
		accessToken, _ := ctx.Cookie("access-token")
		refreshToken, _ := ctx.Cookie("refresh-token")

		claims, newAccessToken, err := r.JWTService.ValidateUserTokens(accessToken, refreshToken)
		if err != nil {
			ctx.AbortWithError(http.StatusForbidden, err)
			return
		}

		// Only requires authentication no permissions
		if len(permissions) == 0 {
			ctx.Set("user", claims)
			if newAccessToken != "" {
				ctx.SetCookie("access-token", newAccessToken, int(r.JWTService.AccessExpires.Seconds()), "/", "", true, true)
			}
			ctx.Next()
			return
		}

		// Check if user has any of the permissions listed
		for _, p := range claims.Permissions {
			if slices.Contains(permissions, p) || p == db.AdminPermission {
				ctx.Set("user", claims)
				if newAccessToken != "" {
					ctx.SetCookie("access-token", newAccessToken, int(r.JWTService.AccessExpires.Seconds()), "/", "", true, true)
				}
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatus(http.StatusForbidden)
	}
}
