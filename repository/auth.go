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
		panic("sas token must be set before the application starts")
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

func (r *AuthRepository) GetProfile(userPrincipal string) (entity.Profile, error) {
	serviceClient := getServiceClient().NewClient("Profiles")
	principalRecord, err := serviceClient.GetEntity(context.TODO(), "actlabs", userPrincipal, nil)
	if err != nil {
		slog.Error("Error getting entity: ", err)
		return entity.Profile{}, err
	}

	profileRecord := entity.ProfileRecord{}
	if err := json.Unmarshal(principalRecord.Value, &profileRecord); err != nil {
		slog.Error("Error unmarshal principal record: ", err)
		return entity.Profile{}, err
	}

	return helper.ConvertRecordToProfile(profileRecord), nil
}

func (r *AuthRepository) GetAllProfiles() ([]entity.Profile, error) {
	profiles := []entity.Profile{}
	profile := entity.Profile{}

	serviceClient := getServiceClient().NewClient("Profiles")

	pager := serviceClient.NewListEntitiesPager(nil)
	for pager.More() {
		response, err := pager.NextPage(context.Background())
		if err != nil {
			slog.Error("Error getting entities: ", err)
			return profiles, err
		}

		for _, entity := range response.Entities {
			var myEntity aztables.EDMEntity
			if err := json.Unmarshal(entity, &myEntity); err != nil {
				slog.Error("Error unmarshal principal record: ", err)
				return profiles, err
			}

			if value, ok := myEntity.Properties["ObjectId"]; ok {
				profile.ObjectId = value.(string)
			} else {
				profile.ObjectId = ""
			}

			if value, ok := myEntity.Properties["DisplayName"]; ok {
				profile.DisplayName = value.(string)
			} else {
				profile.DisplayName = ""
			}

			if value, ok := myEntity.Properties["ProfilePhoto"]; ok {
				profile.ProfilePhoto = value.(string)
			} else {
				profile.ProfilePhoto = ""
			}

			if value, ok := myEntity.Properties["UserPrincipal"]; ok {
				profile.UserPrincipal = value.(string)
			} else {
				profile.UserPrincipal = ""
			}

			if value, ok := myEntity.Properties["Roles"]; ok {
				profile.Roles = helper.StringToSlice(value.(string))
			} else {
				profile.Roles = []string{}
			}

			profiles = append(profiles, profile)
		}
	}

	return profiles, nil
}

// Use this function to complete delete the record for UserPrincipal.
func (r *AuthRepository) DeleteProfile(userPrincipal string) error {
	serviceClient := getServiceClient().NewClient("Profiles")
	_, err := serviceClient.DeleteEntity(context.TODO(), "actlabs", userPrincipal, nil)
	if err != nil {
		slog.Error("Error deleting entity: ", err)
	}
	return err
}

func (r *AuthRepository) UpsertProfile(profile entity.Profile) error {

	//Make sure that profile is complete
	if profile.DisplayName == "" || profile.UserPrincipal == "" {
		slog.Error("Error creating profile: profile is incomplete", nil)
		return errors.New("profile is incomplete")
	}

	serviceClient := getServiceClient().NewClient("Profiles")
	profileRecord := helper.ConvertProfileToRecord(profile)

	marshalledPrincipalRecord, err := json.Marshal(profileRecord)
	if err != nil {
		slog.Error("Error marshalling principal record: ", err)
		return err
	}

	slog.Info("Adding or Updating entity")

	_, err = serviceClient.UpsertEntity(context.TODO(), marshalledPrincipalRecord, nil)
	if err != nil {
		slog.Error("Error adding entity: ", err)
		return err
	}

	return nil
}
