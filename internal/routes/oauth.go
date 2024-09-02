package routes

import (
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/db"
	"github.com/graytonio/flagops-data-storage/internal/db/auth"
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

	dbUser, err := auth.UpsertUser(r.DBClient, db.User{
		Username: user.Name,
		Email: user.Email,
		SSOProvider: user.Provider,
		SSOID: user.UserID,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session := sessions.Default(ctx)
	session.Set("user_id", dbUser.ID)
	if err := session.Save(); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}