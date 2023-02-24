package service

import (
	"actlabs-auth/entity"
	"actlabs-auth/helper"

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

func (s *AuthService) GetRoles(userPrincipal string) ([]string, error) {
	slog.Info("Getting roles for user: ", userPrincipal)
	roles, err := s.authRepository.GetRoles(userPrincipal)
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
	slog.Info("Adding role: ", role, " for user: ", userPrincipal)

	roles, err := s.GetRoles(userPrincipal)
	if err != nil {
		slog.Error("Error getting roles: ", err)
		roles = []string{} // If there is an error, we want to add the role
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
