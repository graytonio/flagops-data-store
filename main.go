package main

import (
	"os"

	"github.com/chenjiandongx/ginprom"
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

	routes := routes.Routes{
		FactProvider: factProvider,
		SecretProvider: secretProvider,
		DBClient: dbClient,
	}

	// TODO Add auth/permissions middleware

	// Managing identities
	r.GET("/identity", routes.GetAllIdentities) // Get all identities
	r.DELETE("/identity/:id", routes.DeleteIdentity) // Delete an identity
	
	// Managing facts
	r.GET("/fact/:id", routes.GetIdentityFacts) // Get all indentity facts
	r.GET("/fact/:id/:fact", routes.GetIdentityFact) // Get specific fact of identity
	r.PUT("/fact/:id/:fact", routes.SetIdentityFact) // Set fact for identity
	r.DELETE("/fact/:id/:fact", routes.DeleteIdentity) // Delete single fact for identity
	
	// Managing secrets
	r.GET("/secret/:id", routes.GetIdentitySecrets) // Get all identity secrets
	r.GET("/secret/:id/:secret", routes.GetIdentitySecret) // Get specific secret of identity
	r.PUT("/secret/:id/:secret", routes.SetIdentitySecret) // Set secret for identity
	r.DELETE("/secret/:id/:secret", routes.DeleteIdentitySecret) // Delete secret for identity

	// Managing users and permissions
	r.GET("/user", routes.RequiresPermission(db.ReadUsers), routes.GetUsers) // Fetch list of users
	r.GET("/user/:id", routes.GetUserByID) // Fetch user details
	r.GET("/permission", routes.GetPermisssions) // Fetch list of available permissions

	r.PUT("/user/:id/permission", routes.AddUserPermissions) // Assign permission to user
	r.DELETE("/user/:id/permission", routes.RemoveUserPermissions) // Remove permission from user

	if err := r.Run(":8080"); err != nil {
		logrus.WithError(err).Fatal("http server crashed")
	}
}