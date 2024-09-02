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
	APIKey      string `gorm:"default:gen_random_uuid();index"`

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
	SecretsWrite    = "secrets-write"
	ReadUsers       = "users-read"
	WriteUsers      = "users-write"
)

var BootstrapPermissions = []Permission{
	{ID: AdminPermission, DisplayName: "Admin"},
	{ID: FactsRead, DisplayName: "Read Facts"},
	{ID: FactsWrite, DisplayName: "Write Facts"},
	{ID: SecretsRead, DisplayName: "Read Secrest"},
	{ID: SecretsWrite, DisplayName: "Write Secrets"},
	{ID: ReadUsers, DisplayName: "Read Users"},
	{ID: WriteUsers, DisplayName: "Write Users"},
}
