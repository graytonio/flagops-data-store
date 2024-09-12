package ui

import (
	"github.com/graytonio/flagops-data-store/internal/config"
	"github.com/graytonio/flagops-data-store/internal/facts"
	"github.com/graytonio/flagops-data-store/internal/secrets"
	"github.com/graytonio/flagops-data-store/internal/services/jwt"
	"github.com/graytonio/flagops-data-store/internal/services/user"
)

// TODO Test routes

type UIRoutes struct {
	Config config.Config

	FactProvider facts.FactProvider
	SecretProvider secrets.SecretProvider

	UserDataService *user.UserDataService
	JWTService *jwt.JWTService
}