package main

import (
	"os"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-data-storage/internal/config"
	"github.com/graytonio/flagops-data-storage/internal/db"
	"github.com/graytonio/flagops-data-storage/internal/routes"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()

	r.Use(routes.ErrorLogger())

	ginPromOpts := ginprom.NewDefaultOpts()
	ginPromOpts.EndpointLabelMappingFn = func(c *gin.Context) string {
		return c.FullPath()
	}

	r.Use(ginprom.PromMiddleware(ginPromOpts))
	r.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	factProvider, secretProvider, err := config.GetProviders()
	if err != nil {
	  logrus.WithError(err).Fatal("could not create providers")
	}

	dbClient, err := db.GetDBClient(os.Getenv("POSTGRES_DB_DSN"))
	if err != nil {
	  logrus.WithError(err).Fatal("could not create db client")
	}

	sessionStore := gormsessions.NewStore(dbClient, true, []byte("secret-salt")) // TODO Generate on first launch or env
	r.Use(sessions.Sessions("cookies", sessionStore))

	routes := routes.Routes{
		FactProvider: factProvider,
		SecretProvider: secretProvider,
		DBClient: dbClient,
	}

	// Managing identities
	r.GET("/identity", routes.RequiresPermission(db.FactsRead, db.SecretsRead), routes.GetAllIdentities) // Get all identities
	r.DELETE("/identity/:id", routes.RequiresPermission(db.FactsRead, db.SecretsRead), routes.DeleteIdentity) // Delete an identity
	
	// Managing facts
	r.GET("/fact/:id", routes.RequiresPermission(db.FactsRead), routes.GetIdentityFacts) // Get all indentity facts
	r.GET("/fact/:id/:fact", routes.RequiresPermission(db.FactsRead), routes.GetIdentityFact) // Get specific fact of identity
	r.PUT("/fact/:id/:fact", routes.RequiresPermission(db.FactsWrite), routes.SetIdentityFact) // Set fact for identity
	r.DELETE("/fact/:id/:fact", routes.RequiresPermission(db.FactsWrite), routes.DeleteIdentity) // Delete single fact for identity
	
	// Managing secrets
	r.GET("/secret/:id", routes.RequiresPermission(db.SecretsRead), routes.GetIdentitySecrets) // Get all identity secrets
	r.GET("/secret/:id/:secret", routes.RequiresPermission(db.SecretsRead), routes.GetIdentitySecret) // Get specific secret of identity
	r.PUT("/secret/:id/:secret", routes.RequiresPermission(db.SecretsWrite), routes.SetIdentitySecret) // Set secret for identity
	r.DELETE("/secret/:id/:secret", routes.RequiresPermission(db.SecretsWrite), routes.DeleteIdentitySecret) // Delete secret for identity

	// Managing users and permissions
	r.GET("/user", routes.RequiresPermission(db.ReadUsers), routes.GetUsers) // Fetch list of users
	r.GET("/user/:id", routes.RequiresPermission(db.ReadUsers), routes.GetUserByID) // Fetch user details
	r.GET("/permission", routes.RequiresPermission(db.ReadUsers), routes.GetPermisssions) // Fetch list of available permissions

	r.PUT("/user/:id/permission", routes.RequiresPermission(db.WriteUsers), routes.AddUserPermissions) // Assign permission to user
	r.DELETE("/user/:id/permission", routes.RequiresPermission(db.WriteUsers), routes.RemoveUserPermissions) // Remove permission from user

	if err := r.Run(":8080"); err != nil {
		logrus.WithError(err).Fatal("http server crashed")
	}
}