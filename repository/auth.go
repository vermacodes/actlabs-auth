package repository

import (
	"actlabs-auth/entity"
	"context"
	"encoding/json"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"golang.org/x/exp/slog"
)

type AuthRepository struct{}

func NewAuthRepository() entity.AuthRepository {
	return &AuthRepository{}
}

func getServiceClient() *aztables.ServiceClient {
	if entity.SasToken == "" || entity.StorageAccountName == "" {
		slog.Error("Error ->", errors.New("sas token and storage account name must be set before the application starts"))
		panic("sas token and storage account name must be set before the application starts")
	}
	SasUrl := "https://" + entity.StorageAccountName + ".table.core.windows.net/" + entity.SasToken

	serviceClient, err := aztables.NewServiceClientFromConnectionString(SasUrl, nil)
	if err != nil {
		slog.Error("Error creating service client: ", err)
	}

	return serviceClient
}

func (r *AuthRepository) GetRoles(userPrincipal string) ([]string, error) {
	roles := []string{}
	serviceClient := getServiceClient().NewClient("Roles")
	principalRecord, err := serviceClient.GetEntity(context.Background(), "actlabs", userPrincipal, nil)
	if err != nil {
		slog.Error("Error getting entity: ", err)
		return roles, err
	}

	roleRecord := entity.RoleRecord{}
	if err := json.Unmarshal(principalRecord.Value, &roleRecord); err != nil {
		slog.Error("Error unmarshalling principal record: ", err)
		return roles, err
	}

	return roleRecord.Roles, nil
}

// Use this function to complete delete the record for UserPricipal.
func (r *AuthRepository) DeleteRole(userPrincipal string) error {
	serviceClient := getServiceClient().NewClient("Roles")
	_, err := serviceClient.DeleteEntity(context.Background(), "actlabs", userPrincipal, nil)
	if err != nil {
		slog.Error("Error deleting entity: ", err)
	}
	return err
}

func (r *AuthRepository) AddRole(userPrincipal string, roles []string) error {
	serviceClient := getServiceClient().NewClient("Roles")
	principalRecord := entity.RoleRecord{
		Entity: aztables.Entity{
			PartitionKey: "actlabs",
			RowKey:       userPrincipal,
		},
		UserPrincipal: userPrincipal,
		Roles:         roles,
	}

	marshalledPrincipalRecord, err := json.Marshal(principalRecord)
	if err != nil {
		slog.Error("Error marshalling principal record: ", err)
		return err
	}

	_, err = serviceClient.AddEntity(context.Background(), marshalledPrincipalRecord, nil)
	if err != nil {
		slog.Error("Error adding entity: ", err)
		return err
	}

	return nil
}
