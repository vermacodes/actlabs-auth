package service

import (
	"actlabs-auth/entity"

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
	slog.Debug("Getting roles for user: ", userPrincipal)
	roles, err := s.authRepository.GetRoles(userPrincipal)
	if err != nil {
		slog.Error("Error getting roles: ", err)
	}
	return roles, err
}

func (s *AuthService) DeleteRole(userPrincipal string, role string) error {
	slog.Debug("Deleting role: ", role, " for user: ", userPrincipal)

	roles, err := s.authRepository.GetRoles(userPrincipal)
	if err != nil {
		slog.Error("Error getting roles: ", err)
		return err
	}
	remove(roles, role)
	if len(roles) == 0 {
		return s.authRepository.DeleteRole(userPrincipal)
	}

	// This adds the roles again after removing the role
	return s.authRepository.AddRole(userPrincipal, roles)
}

func (s *AuthService) AddRole(userPrincipal string, role string) error {
	slog.Debug("Adding role: ", role, " for user: ", userPrincipal)

	roles, err := s.authRepository.GetRoles(userPrincipal)
	if err != nil {
		slog.Error("Error getting roles: ", err)
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
