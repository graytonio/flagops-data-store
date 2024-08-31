package routes

import (
	"errors"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/auth"
	"github.com/graytonio/flagops-data-storage/internal/db"
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

// Fetches user either from session info or from auth token
func (r *Routes) getSessionUser(ctx *gin.Context) (*db.User, error) {
	// TODO Attempt to fetch from session

	rawUserID := ctx.Request.Header.Get("Authorization") // TODO Make real api token
	if rawUserID == "" {
		return nil, errors.New("id parameter must not be empty")
	}

	userId, err := strconv.ParseUint(rawUserID, 10, 0)
	if err != nil {
		return nil, err
	}

	user, err := auth.GetUserByID(r.DBClient, uint(userId))
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Path protected to require a user session
func (r *Routes) RequiresAuthentication(ctx *gin.Context) {
	_, err := r.getSessionUser(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	ctx.Next()
}

// Path protected by specific user permissions
func (r *Routes) RequiresPermission(permissions ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := r.getSessionUser(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		for _, p := range user.Permissions {
			if p.ID == db.AdminPermission {
				ctx.Next()
				return
			}

			if slices.Contains(permissions, p.ID) {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatus(http.StatusForbidden)
	}
}
