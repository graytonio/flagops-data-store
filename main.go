package main

import (
	"net/http"
	"time"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/config"
	"github.com/graytonio/flagops-data-storage/internal/db"
	"github.com/graytonio/flagops-data-storage/internal/facts"
	"github.com/graytonio/flagops-data-storage/internal/renderer"
	"github.com/graytonio/flagops-data-storage/internal/routes"
	"github.com/graytonio/flagops-data-storage/internal/routes/api"
	"github.com/graytonio/flagops-data-storage/internal/routes/ui"
	"github.com/graytonio/flagops-data-storage/internal/secrets"
	"github.com/graytonio/flagops-data-storage/internal/services/jwt"
	"github.com/graytonio/flagops-data-storage/internal/services/user"
	"github.com/graytonio/flagops-data-storage/templates/pages"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	conf, err := config.ParseConfig()
	if err != nil {
		logrus.WithError(err).Fatal("cannot parse config")
	}

	dbClient, err := db.GetDBClient(conf.UserDatabaseOptions.PostgresDSN)
	if err != nil {
		logrus.WithError(err).Fatal("could not connect to db deployment")
	}

	factProvider, err := facts.GetFactProvider(conf.FactsProviderOptions)
	if err != nil {
		logrus.WithError(err).Fatal("cannot init fact provider")
	}

	secretProvider, err := secrets.GetSecretsProvider(conf.SecretsProviderOptions)
	if err != nil {
		logrus.WithError(err).Fatal("cannot init secrets provider")
	}

	userDataService := &user.UserDataService{
		DBClient: dbClient,
	}

	jwtService := &jwt.JWTService{
		AccessExpires:   time.Minute * time.Duration(conf.UserDatabaseOptions.AccessTokenExpirationMinutes),
		RefreshExpires:  time.Minute * time.Duration(conf.UserDatabaseOptions.RefreshTokenExpirationMinutes),
		SigningSecret:   conf.UserDatabaseOptions.JWTSecret,
		UserDataService: userDataService,
	}

	routeHandlers := &routes.Routes{
		Config: *conf,

		FactProvider:   factProvider,
		SecretProvider: secretProvider,

		UserDataService: userDataService,
		JWTService:      jwtService,
	}
	routeHandlers.InitOauthProvider()

	apiRoutesHandlers := &api.APIRoutes{
		Config: *conf,

		FactProvider:   factProvider,
		SecretProvider: secretProvider,

		UserDataService: userDataService,
		JWTService:      jwtService,
	}

	uiRoutesHandlers := &ui.UIRoutes{
		Config: *conf,

		FactProvider:   factProvider,
		SecretProvider: secretProvider,

		UserDataService: userDataService,
		JWTService:      jwtService,
	}

	r := gin.Default()

	r.HTMLRender = &renderer.HTMLTemplRenderer{}
	r.Static("/assets", "/assets")

	r.Use(routes.ErrorLogger())

	ginPromOpts := ginprom.NewDefaultOpts()
	ginPromOpts.EndpointLabelMappingFn = func(c *gin.Context) string {
		return c.FullPath()
	}

	r.Use(ginprom.PromMiddleware(ginPromOpts))
	r.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	apiRoutes := r.Group("/api")
	{
		// Managing identities
		apiRoutes.GET("/identity", /*routeHandlers.RequiresAuth(db.FactsRead, db.SecretsRead),*/ apiRoutesHandlers.GetAllIdentities)      // Get all identities
		apiRoutes.DELETE("/identity/:id", routeHandlers.RequiresAuth(db.FactsRead, db.SecretsRead), apiRoutesHandlers.DeleteIdentity) // Delete an identity

		// Managing facts
		apiRoutes.GET("/fact/:id", /*routeHandlers.RequiresAuth(db.FactsRead),*/ apiRoutesHandlers.GetIdentityFacts)         // Get all indentity facts
		apiRoutes.GET("/fact/:id/:fact", routeHandlers.RequiresAuth(db.FactsRead), apiRoutesHandlers.GetIdentityFact)    // Get specific fact of identity
		apiRoutes.PUT("/fact/:id/:fact", /*routeHandlers.RequiresAuth(db.FactsWrite),*/ apiRoutesHandlers.SetIdentityFact)   // Set fact for identity
		apiRoutes.DELETE("/fact/:id/:fact", routeHandlers.RequiresAuth(db.FactsWrite), apiRoutesHandlers.DeleteIdentity) // Delete single fact for identity

		// Managing secrets
		apiRoutes.GET("/secret/:id", routeHandlers.RequiresAuth(db.SecretsRead), apiRoutesHandlers.GetIdentitySecrets)               // Get all identity secrets
		apiRoutes.GET("/secret/:id/:secret", routeHandlers.RequiresAuth(db.SecretsRead), apiRoutesHandlers.GetIdentitySecret)        // Get specific secret of identity
		apiRoutes.PUT("/secret/:id/:secret", routeHandlers.RequiresAuth(db.SecretsWrite), apiRoutesHandlers.SetIdentitySecret)       // Set secret for identity
		apiRoutes.DELETE("/secret/:id/:secret", routeHandlers.RequiresAuth(db.SecretsWrite), apiRoutesHandlers.DeleteIdentitySecret) // Delete secret for identity

		// Managing users and permissions
		apiRoutes.GET("/user", routeHandlers.RequiresAuth(db.ReadUsers), apiRoutesHandlers.GetUsers)                                 // Fetch list of users
		apiRoutes.GET("/user/:id", routeHandlers.RequiresAuth(db.ReadUsers), apiRoutesHandlers.GetUserByID)                          // Fetch user details
		apiRoutes.GET("/permission", routeHandlers.RequiresAuth(db.ReadUsers), apiRoutesHandlers.GetPermisssions)                    // Fetch list of available permissions
		apiRoutes.PUT("/user/:id/permission", routeHandlers.RequiresAuth(db.WriteUsers), apiRoutesHandlers.AddUserPermissions)       // Assign permission to user
		apiRoutes.DELETE("/user/:id/permission", routeHandlers.RequiresAuth(db.WriteUsers), apiRoutesHandlers.RemoveUserPermissions) // Remove permission from user
	}

	uiRoutes := r.Group("/ui")
	{
		uiRoutes.GET("/", uiRoutesHandlers.HomeDashboard)
		uiRoutes.GET("/fact/:id", uiRoutesHandlers.IdentityFactsDashboard)
		uiRoutes.GET("/htmx/fact/:id", uiRoutesHandlers.IdentityFactsTable)
	}

	// Authentication
	r.GET("/auth/login", routeHandlers.OauthLogin)
	r.GET("/auth/github/callback", routeHandlers.OauthCallback)

	r.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "", pages.LoginPage(conf.OAuthOptions.Provider))
	})

	if err := r.Run(":8080"); err != nil {
		logrus.WithError(err).Fatal("http server crashed")
	}
}
