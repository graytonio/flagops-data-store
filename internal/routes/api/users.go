package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func (r *APIRoutes) GetUsers(ctx *gin.Context) {
	users, err := r.UserDataService.GetUsers()
	if err != nil {
	  ctx.AbortWithError(http.StatusInternalServerError, err)
	  return
	}

	ctx.JSON(http.StatusOK, users)
}

func (r *APIRoutes) GetUserByID(ctx *gin.Context) {
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

	user, err := r.UserDataService.GetUserByID(uint(userId))
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

func (r *APIRoutes) AddUserPermissions(ctx *gin.Context) {
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

	err = r.UserDataService.AddUserPermissions(uint(userId), body.Permissions)
	if err != nil {
	  ctx.AbortWithError(http.StatusInternalServerError, err)
	  return
	}
}

func (r *APIRoutes) RemoveUserPermissions(ctx *gin.Context) {
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

	err = r.UserDataService.RemoveUserPermissions(uint(userId), body.Permissions)
	if err != nil {
	  ctx.AbortWithError(http.StatusInternalServerError, err)
	  return
	}
}

func (r *APIRoutes) GetPermisssions(ctx *gin.Context) {
	permissions, err := r.UserDataService.GetPermissions()
	if err != nil {
	  ctx.AbortWithError(http.StatusInternalServerError, err)
	  return
	}

	ctx.JSON(http.StatusOK, permissions)
}