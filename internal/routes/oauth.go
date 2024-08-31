package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/auth"
	"github.com/graytonio/flagops-data-storage/internal/db"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/oov/gothic"
)

func init() {
	goth.UseProviders(
		github.New(os.Getenv("GITHUB_OAUTH_CLIENT_KEY"), os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"), os.Getenv("HOSTNAME")+"/auth/github/callback"),
	)
}

func (r *Routes) OauthLogin(ctx *gin.Context) {
	err := gothic.BeginAuth(os.Getenv("OAUTH_PROVIDER"), ctx.Writer, ctx.Request)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (r *Routes) OauthCallback(ctx *gin.Context) {
	user, err := gothic.CompleteAuth(os.Getenv("OAUTH_PROVIDER"), ctx.Writer, ctx.Request)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	_, err = auth.CreateOrUpdateUser(r.DBClient, db.User{
		Username: user.Name,
		Email: user.Email,
		SSOProvider: user.Provider,
		SSOID: user.UserID,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// TODO Store JWT in session

	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}