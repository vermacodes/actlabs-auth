package service

import (
	"actlabs-auth/entity"
	"actlabs-auth/helper"
	"strings"

	"golang.org/x/exp/slog"
)

type AuthService struct {
	authRepository entity.AuthRepository
}

func NewAuthService(authRepository entity.AuthRepository) entity.AuthService {
	return &AuthService{
		authRepository: authRepository,
	}
}

func (s *AuthService) GetRoles(userPrincipal string) (entity.Roles, error) {
	slog.Info("Getting roles for user: ", userPrincipal)
	roles, err := s.authRepository.GetRoles(userPrincipal)
	if err != nil {
		slog.Error("Error getting roles: ", err)

		// if erro contains "ERROR CODE: ResourceNotFound" then add the user role
		if strings.Contains(err.Error(), "ERROR CODE: ResourceNotFound") {
			roles := []string{"user"}
			if err := s.authRepository.AddRole(userPrincipal, roles); err != nil {
				slog.Error("Error adding role: ", err)
			}
		}
	}

	// Add the roles to the Roles object.
	rolesObj := entity.Roles{
		UserPrincipal: userPrincipal,
		Roles:         roles,
	}
	return rolesObj, err
}

func (s *AuthService) GetAllRoles() ([]entity.Roles, error) {
	slog.Info("Getting all roles")
	roles, err := s.authRepository.GetAllRoles()
	if err != nil {
		slog.Error("Error getting roles: ", err)
	}
	return roles, err
}

func (s *AuthService) DeleteRole(userPrincipal string, role string) error {
	slog.Info("Deleting role: ", role, " for user: ", userPrincipal)

	roles, err := s.authRepository.GetRoles(userPrincipal)
	if err != nil {
		slog.Error("Error getting roles: ", err)
		return err
	}
	roles = remove(roles, role)
	if len(roles) == 0 {
		return s.authRepository.DeleteRole(userPrincipal)
	}

	// This adds the roles again after removing the role
	return s.authRepository.AddRole(userPrincipal, roles)
}

func (s *AuthService) AddRole(userPrincipal string, role string) error {
	slog.Info("Adding role: " + role + " for user: " + userPrincipal)

	roles := []string{}

	rolesObj, err := s.GetRoles(userPrincipal)
	if err != nil {
		slog.Error("Error getting roles: ", err)
	} else {
		roles = rolesObj.Roles
	}

	if helper.Contains(roles, role) {
		slog.Info("Role already exists: ", role)
		return nil
	}

	roles = append(roles, role)

	return s.authRepository.AddRole(userPrincipal, roles)
}

// Helper Function to remove an element from a slice
func remove(roles []string, role string) []string {
	for i, v := range roles {
		if v == role {
			roles = append(roles[:i], roles[i+1:]...)
			break
		}
	}
	return roles
}
