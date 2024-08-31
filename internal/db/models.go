package db

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string

	SSOProvider string
	SSOID       string
	APIKey      string

	Permissions []Permission `gorm:"many2many:user_permissions"`
}

type Permission struct {
	ID          string `gorm:"uniqueIndex"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	DisplayName string
}

const (
	AdminPermission = "admin"
	FactsRead       = "facts-read"
	FactsWrite      = "facts-write"
	SecretsRead     = "secrets-read"
	SecretsWrite    = "secrest-write"
	ReadUsers       = "users-read"
	WriteUsers      = "users-write"
)

var BootstrapPermissions = []Permission{
	{ID: AdminPermission},
	{ID: FactsRead},
	{ID: FactsWrite},
	{ID: SecretsRead},
	{ID: SecretsWrite},
	{ID: ReadUsers},
	{ID: WriteUsers},
}
