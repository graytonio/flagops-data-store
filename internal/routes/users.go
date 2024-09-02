package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/db/auth"
	"gorm.io/gorm"
)


func (r *Routes) GetUsers(ctx *gin.Context) {
	users, err := auth.GetUsers(r.DBClient)
	if err != nil {
	  ctx.AbortWithError(http.StatusInternalServerError, err)
	  return
	}

	ctx.JSON(http.StatusOK, users)
}

func (r *Routes) GetUserByID(ctx *gin.Context) {
	rawUserID := ctx.Param("id")
	if rawUserID == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter must not be empty"))
		return
	}

	userId, err := strconv.ParseUint(rawUserID, 10, 0)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid user id"))
		return
	}

	user, err := auth.GetUserByID(r.DBClient, uint(userId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type modifyPermissionRequest struct {
	Permissions []string `json:"permissions"`
}

func (r *Routes) AddUserPermissions(ctx *gin.Context) {
	rawUserID := ctx.Param("id")
	if rawUserID == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter must not be empty"))
		return
	}

	userId, err := strconv.ParseUint(rawUserID, 10, 0)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid user id"))
		return
	}
	
	var body modifyPermissionRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = auth.AddUserPermissions(r.DBClient, uint(userId), body.Permissions)
	if err != nil {
	  ctx.AbortWithError(http.StatusInternalServerError, err)
	  return
	}
}

func (r *Routes) RemoveUserPermissions(ctx *gin.Context) {
	rawUserID := ctx.Param("id")
	if rawUserID == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("id parameter must not be empty"))
		return
	}

	userId, err := strconv.ParseUint(rawUserID, 10, 0)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid user id"))
		return
	}
	
	var body modifyPermissionRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = auth.RemoveUserPermissions(r.DBClient, uint(userId), body.Permissions)
	if err != nil {
	  ctx.AbortWithError(http.StatusInternalServerError, err)
	  return
	}
}

func (r *Routes) GetPermisssions(ctx *gin.Context) {
	permissions, err := auth.GetPermissions(r.DBClient)
	if err != nil {
	  ctx.AbortWithError(http.StatusInternalServerError, err)
	  return
	}

	ctx.JSON(http.StatusOK, permissions)
}