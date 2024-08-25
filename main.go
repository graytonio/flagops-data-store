package main

import (
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/graytonio/flagops-config-storage/internal/config"
	"github.com/graytonio/flagops-config-storage/internal/routes"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()

	r.Use(ginprom.PromMiddleware(ginprom.NewDefaultOpts()))
	r.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))


	factProvider, secretProvider, err := config.GetProviders()
	if err != nil {
	  logrus.WithError(err).Fatal("could not create providers")
	}

	routes := routes.Routes{
		FactProvider: factProvider,
		SecretProvider: secretProvider,
	}

	// Managing identities
	r.GET("/identity", routes.GetAllIdentities) // Get all identities
	r.DELETE("/identity/:id", routes.DeleteIdentity) // Delete an identity
	
	// Managing facts
	r.GET("/fact/:id", routes.GetIdentityFacts) // Get all indentity facts
	r.GET("/fact/:id/:fact", routes.GetIdentityFact) // Get specific fact of identity
	r.PUT("/fact/:id/:fact", routes.SetIdentityFact) // Set fact for identity
	r.DELETE("/fact/:id/:fact", routes.DeleteIdentity) // Delete single fact for identity
	
	// Managing secrets
	r.GET("/secrets/:id", routes.GetIdentitySecrets) // Get all identity secrets
	r.GET("/fact/:id/:secret", routes.GetIdentitySecret) // Get specific secret of identity
	r.PUT("/fact/:id/:secret", routes.SetIdentitySecret) // Set secret for identity
	r.DELETE("/secret/:id/:secret", routes.DeleteIdentitySecret) // Delete secret for identity

	if err := r.Run(":8080"); err != nil {
		logrus.WithError(err).Fatal("http server crashed")
	}
}