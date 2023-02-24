package repository

import (
	"actlabs-auth/entity"
	"actlabs-auth/helper"
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
	if entity.SasToken == "" {
		slog.Error("Error ->", errors.New("sas token must be set before the application starts"))
		panic("sas tokenmust be set before the application starts")
	}

	if entity.StorageAccountName == "" {
		slog.Error("Error ->", errors.New("storage account name must be set before the application starts"))
		panic("storage account name must be set before the application starts")
	}

	SasUrl := "https://" + entity.StorageAccountName + ".table.core.windows.net/" + entity.SasToken

	serviceClient, err := aztables.NewServiceClientWithNoCredential(SasUrl, nil)
	if err != nil {
		slog.Error("Error creating service client: ", err)
	}

	return serviceClient
}

func (r *AuthRepository) GetRoles(userPrincipal string) ([]string, error) {
	roles := []string{}
	serviceClient := getServiceClient().NewClient("Roles")
	principalRecord, err := serviceClient.GetEntity(context.TODO(), "actlabs", userPrincipal, nil)
	if err != nil {
		slog.Error("Error getting entity: ", err)
		return roles, err
	}

	slog.Info("principalRecord: ", string(principalRecord.Value))

	roleRecord := entity.RoleRecord{}
	if err := json.Unmarshal(principalRecord.Value, &roleRecord); err != nil {
		slog.Error("Error unmarshalling principal record: ", err)
		return roles, err
	}

	return helper.StringToSlice(roleRecord.Roles), nil
}

// Use this function to complete delete the record for UserPricipal.
func (r *AuthRepository) DeleteRole(userPrincipal string) error {
	serviceClient := getServiceClient().NewClient("Roles")
	_, err := serviceClient.DeleteEntity(context.TODO(), "actlabs", userPrincipal, nil)
	if err != nil {
		slog.Error("Error deleting entity: ", err)
	}
	return err
}

func (r *AuthRepository) AddRole(userPrincipal string, roles []string) error {
	serviceClient := getServiceClient().NewClient("Roles")
	principalRecord := entity.RoleRecord{
		PartitionKey:  "actlabs",
		RowKey:        userPrincipal,
		UserPrincipal: userPrincipal,
		Roles:         helper.SliceToString(roles),
	}

	marshalledPrincipalRecord, err := json.Marshal(principalRecord)
	if err != nil {
		slog.Error("Error marshalling principal record: ", err)
		return err
	}

	slog.Info("Adding entity: " + string(marshalledPrincipalRecord))

	_, err = serviceClient.UpsertEntity(context.TODO(), marshalledPrincipalRecord, nil)
	if err != nil {
		slog.Error("Error adding entity: ", err)
		return err
	}

	return nil
}
