package routes

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/db"
	"github.com/graytonio/flagops-data-storage/internal/db/auth"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/oov/gothic"
	"github.com/spf13/viper"
)

func init() {
	goth.UseProviders(
		github.New(viper.GetString("GITHUB_OAUTH_CLIENT_KEY"), viper.GetString("GITHUB_OAUTH_CLIENT_SECRET"), viper.GetString("HOSTNAME")+"/auth/github/callback"),
	)
}

func (r *Routes) OauthLogin(ctx *gin.Context) {
	err := gothic.BeginAuth(viper.GetString("OAUTH_PROVIDER"), ctx.Writer, ctx.Request)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (r *Routes) OauthCallback(ctx *gin.Context) {
	user, err := gothic.CompleteAuth(viper.GetString("OAUTH_PROVIDER"), ctx.Writer, ctx.Request)
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