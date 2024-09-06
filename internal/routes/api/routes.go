package api

import (
	"github.com/graytonio/flagops-data-storage/internal/config"
	"github.com/graytonio/flagops-data-storage/internal/facts"
	"github.com/graytonio/flagops-data-storage/internal/secrets"
	"github.com/graytonio/flagops-data-storage/internal/services/jwt"
	"github.com/graytonio/flagops-data-storage/internal/services/user"
)

// TODO Test routes

type APIRoutes struct {
	Config config.Config

	FactProvider facts.FactProvider
	SecretProvider secrets.SecretProvider

	UserDataService *user.UserDataService
	JWTService *jwt.JWTService
}
