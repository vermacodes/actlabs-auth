package entity

import (
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/gin-gonic/gin"
)

// These variables are used to store the Azure Storage Account name and SAS token.
// They must be set before the application starts.
var SasToken string
var StorageAccountName string

type Auth struct {
	UserPrincipal string   `json:"userPrincipal"`
	Roles         []string `json:"roles"`
}

type RoleRecord struct {
	Entity        aztables.Entity
	UserPrincipal string
	Roles         []string
}

type AuthService interface {
	// Get Roles
	GetRoles(userPrincipal string) ([]string, error)
	DeleteRole(userPrincipal string, role string) error
	AddRole(userPrincipal string, role string) error
}

type AuthHandler interface {
	// Get Roles
	GetRoles(c *gin.Context)
	DeleteRole(c *gin.Context)
	AddRole(c *gin.Context)
}

type AuthRepository interface {
	// Get Roles
	GetRoles(userPrincipal string) ([]string, error)

	// This method is used to delete the record for UserPricipal from the table.
	// This is used only when the last role is removed from the user.
	DeleteRole(userPrincipal string) error
	AddRole(userPrincipal string, roles []string) error
}
