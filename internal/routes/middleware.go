package routes

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/db"
	"github.com/graytonio/flagops-data-storage/internal/db/auth"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
func (r *Routes) getAuthUser(ctx *gin.Context) (*db.User, error) {
	user, err := r.getSessionUser(ctx)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	user, err = r.getAPIUser(ctx)
	if err != nil {
	  return nil, err
	}

	if user != nil {
		return user, nil
	}

	return nil, nil
}

func (r *Routes) getAPIUser(ctx *gin.Context) (*db.User, error) {
	apiKey := strings.TrimPrefix(ctx.Request.Header.Get("Authorization"), "Bearer ")
	if apiKey == "" {
		return nil, nil
	}

	user, err := auth.GetUserByAPIKey(r.DBClient, apiKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r *Routes) getSessionUser(ctx *gin.Context) (*db.User, error) {
	session := sessions.Default(ctx)
	rawUserID := session.Get("user_id")
	if rawUserID == nil {
		return nil, nil
	}

	userId, err := strconv.ParseUint(rawUserID.(string), 10, 0)
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
	user, err := r.getAuthUser(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	ctx.Next()
}

// A route that can only be called by the user themselves or an admin
func (r *Routes) RequiresSelf(ctx *gin.Context) {
	user, err := r.getAuthUser(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Allow admins
	for _, p := range user.Permissions {
		if p.ID == db.AdminPermission {
			ctx.Next()
			return
		}
	}

	userId, err := strconv.ParseUint(ctx.Param("id"), 10, 0)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if user.ID != uint(userId) {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	ctx.Next()
}

// Path protected by specific user permissions
func (r *Routes) RequiresPermission(permissions ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := r.getAuthUser(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	
		if user == nil {
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
