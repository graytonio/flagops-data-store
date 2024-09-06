package main

import (
	"time"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/config"
	"github.com/graytonio/flagops-data-storage/internal/db"
	"github.com/graytonio/flagops-data-storage/internal/facts"
	"github.com/graytonio/flagops-data-storage/internal/routes"
	"github.com/graytonio/flagops-data-storage/internal/secrets"
	"github.com/graytonio/flagops-data-storage/internal/services/jwt"
	"github.com/graytonio/flagops-data-storage/internal/services/user"
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
		AccessExpires: time.Minute * time.Duration(conf.UserDatabaseOptions.AccessTokenExpirationMinutes),
		RefreshExpires: time.Minute * time.Duration(conf.UserDatabaseOptions.RefreshTokenExpirationMinutes),
		SigningSecret: conf.UserDatabaseOptions.JWTSecret,
		UserDataService: userDataService,
	}

	routeHandlers := &routes.Routes{
		Config: *conf,

		FactProvider: factProvider,
		SecretProvider: secretProvider,

		UserDataService: userDataService,
		JWTService: jwtService,
	}
	routeHandlers.InitOauthProvider()

	r := gin.Default()

	r.Use(routes.ErrorLogger())
	
	ginPromOpts := ginprom.NewDefaultOpts()
	ginPromOpts.EndpointLabelMappingFn = func(c *gin.Context) string {
		return c.FullPath()
	}

	r.Use(ginprom.PromMiddleware(ginPromOpts))
	r.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	// Managing identities
	r.GET("/identity", routeHandlers.RequiresAuth(db.FactsRead, db.SecretsRead), routeHandlers.GetAllIdentities) // Get all identities
	r.DELETE("/identity/:id", routeHandlers.RequiresAuth(db.FactsRead, db.SecretsRead), routeHandlers.DeleteIdentity) // Delete an identity
	
	// Managing facts
	r.GET("/fact/:id", routeHandlers.RequiresAuth(db.FactsRead), routeHandlers.GetIdentityFacts) // Get all indentity facts
	r.GET("/fact/:id/:fact", routeHandlers.RequiresAuth(db.FactsRead), routeHandlers.GetIdentityFact) // Get specific fact of identity
	r.PUT("/fact/:id/:fact", routeHandlers.RequiresAuth(db.FactsWrite), routeHandlers.SetIdentityFact) // Set fact for identity
	r.DELETE("/fact/:id/:fact", routeHandlers.RequiresAuth(db.FactsWrite), routeHandlers.DeleteIdentity) // Delete single fact for identity
	
	// Managing secrets
	r.GET("/secret/:id", routeHandlers.RequiresAuth(db.SecretsRead), routeHandlers.GetIdentitySecrets) // Get all identity secrets
	r.GET("/secret/:id/:secret", routeHandlers.RequiresAuth(db.SecretsRead), routeHandlers.GetIdentitySecret) // Get specific secret of identity
	r.PUT("/secret/:id/:secret", routeHandlers.RequiresAuth(db.SecretsWrite), routeHandlers.SetIdentitySecret) // Set secret for identity
	r.DELETE("/secret/:id/:secret", routeHandlers.RequiresAuth(db.SecretsWrite), routeHandlers.DeleteIdentitySecret) // Delete secret for identity

	// Managing users and permissions
	r.GET("/user", routeHandlers.RequiresAuth(db.ReadUsers), routeHandlers.GetUsers) // Fetch list of users
	r.GET("/user/:id", routeHandlers.RequiresAuth(db.ReadUsers), routeHandlers.GetUserByID) // Fetch user details
	r.GET("/permission", routeHandlers.RequiresAuth(db.ReadUsers), routeHandlers.GetPermisssions) // Fetch list of available permissions
	r.PUT("/user/:id/permission", routeHandlers.RequiresAuth(db.WriteUsers), routeHandlers.AddUserPermissions) // Assign permission to user
	r.DELETE("/user/:id/permission", routeHandlers.RequiresAuth(db.WriteUsers), routeHandlers.RemoveUserPermissions) // Remove permission from user

	// Authentication
	r.GET("/auth/login", routeHandlers.OauthLogin)
	r.GET("/auth/github/callback", routeHandlers.OauthCallback)

	if err := r.Run(":8080"); err != nil {
		logrus.WithError(err).Fatal("http server crashed")
	}
}