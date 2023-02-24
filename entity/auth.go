package entity

import (
	"github.com/gin-gonic/gin"
)

// These variables are used to store the Azure Storage Account name and SAS token.
// They must be set before the application starts.
var SasToken string
var StorageAccountName string

type UserRole struct {
	UserPrincipal string   `json:"userPrincipal"`
	Roles         []string `json:"roles"`
}

// For some reason, I am not able to add Roles as a slice of strings to the table.
// So, I am converting the slice to a string and then converting it back to a slice.
// This is not the best way to do it, but it works for now.
type RoleRecord struct {
	PartitionKey  string `json:"PartitionKey"`
	RowKey        string `json:"RowKey"`
	UserPrincipal string `json:"UserPrincipal"`
	Roles         string `json:"Roles"`
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
