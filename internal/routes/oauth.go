package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/db"
	"github.com/graytonio/flagops-data-storage/internal/services/jwt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/oov/gothic"
	"github.com/sirupsen/logrus"
)

func (r *Routes) InitOauthProvider() {
	var provider goth.Provider
	switch r.Config.OAuthOptions.Provider {
	case "github":
		provider = github.New(
			r.Config.OAuthOptions.GithubClientKey,
			r.Config.OAuthOptions.GithubClientSecret,
			r.Config.OAuthOptions.Hostname + "/auth/github/callback",
		)
	default:
		logrus.WithField("oauth_provider", r.Config.OAuthOptions.Provider).Debug("unsupported oauth provider")
		return
	}
	goth.UseProviders(provider)
}

func (r *Routes) OauthLogin(ctx *gin.Context) {
	err := gothic.BeginAuth(r.Config.OAuthOptions.Provider, ctx.Writer, ctx.Request)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (r *Routes) OauthCallback(ctx *gin.Context) {
	user, err := gothic.CompleteAuth(r.Config.OAuthOptions.Provider, ctx.Writer, ctx.Request)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	dbUser, err := r.UserDataService.UpsertUser(db.User{
		Username: user.Name,
		Email: user.Email,
		SSOProvider: user.Provider,
		SSOID: user.UserID,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	permissions := []string{}
	for _, p := range dbUser.Permissions {
		permissions = append(permissions, p.ID)
	}

	accessToken, err := r.JWTService.NewUserAccessToken(&jwt.UserClaims{
		ID: dbUser.ID,
		Permissions: permissions,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	refreshToken, err := r.JWTService.NewUserRefreshToken(&jwt.UserRefreshClaims{
		ID: dbUser.ID,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.SetCookie("access-token", accessToken, int(r.JWTService.AccessExpires.Seconds()), "/", "", true, true)
	ctx.SetCookie("refresh-token", refreshToken, int(r.JWTService.RefreshExpires.Seconds()), "/", "", true, true)

	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}