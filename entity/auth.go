package entity

import (
	"github.com/gin-gonic/gin"
)

// These variables are used to store the Azure Storage Account name and SAS token.
// They must be set before the application starts.
var SasToken string
var StorageAccountName string

type Profile struct {
	ObjectId        string   `json:"objectId"`
	UserPrincipal   string   `json:"userPrincipal"`
	DisplayName     string   `json:"displayName"`
	ProfilePhotoUrl string   `json:"profilePhotoUrl"`
	Roles           []string `json:"roles"`
}

// Azure storage table doesn't support adding an array of strings. Thus, the hack.
// This is not the best way to do it, but it works for now.
type ProfileRecord struct {
	PartitionKey    string `json:"PartitionKey"`
	RowKey          string `json:"RowKey"`
	ObjectId        string `json:"ObjectId"`
	UserPrincipal   string `json:"UserPrincipal"`
	DisplayName     string `json:"DisplayName"`
	ProfilePhotoUrl string `json:"ProfilePhotoUrl"`
	Roles           string `json:"Roles"`
}

type AuthService interface {
	// Get Profile
	GetProfile(userPrincipal string) (Profile, error)
	GetAllProfiles() ([]Profile, error)
	DeleteRole(userPrincipal string, role string) error
	AddRole(userPrincipal string, role string) error
}

type AuthHandler interface {
	// Get Roles
	GetProfile(c *gin.Context)
	GetAllProfiles(c *gin.Context)
	DeleteRole(c *gin.Context)
	AddRole(c *gin.Context)
}

type AuthRepository interface {
	// Get Roles
	GetProfile(userPrincipal string) (Profile, error)
	GetAllProfiles() ([]Profile, error)

	// This method is used to delete the record for UserPrincipal from the table.
	// This is used only when the last role is removed from the user.
	DeleteProfile(userPrincipal string) error

	// This method is used to add a role to the user.
	UpsertProfile(profile Profile) error
}
