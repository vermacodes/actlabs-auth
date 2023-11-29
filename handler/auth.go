package handler

import (
	"actlabs-auth/entity"
	"actlabs-auth/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type AuthHandler struct {
	authService entity.AuthService
}

func NewAuthHandler(r *gin.RouterGroup, authService entity.AuthService) {
	handler := &AuthHandler{
		authService: authService,
	}

	r.GET("/profiles/:userPrincipal", handler.GetProfile)
	r.POST("/profiles/default", handler.AddDefaultProfile)
}

func NewAdminAuthHandler(r *gin.RouterGroup, authService entity.AuthService) {
	handler := &AuthHandler{
		authService: authService,
	}

	r.GET("/profiles", handler.GetAllProfiles)
	r.POST("/profiles/:userPrincipal/:role", handler.AddRole)
	r.DELETE("/profiles/:userPrincipal/:role", handler.DeleteRole)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userPrincipal := c.Param("userPrincipal")

	// My roles
	if userPrincipal == "my" {

		// Get the auth token from the request header
		authToken := c.GetHeader("Authorization")

		// Remove Bearer from the authToken
		authToken = strings.Split(authToken, "Bearer ")[1]

		userPrincipal, _ = helper.GetUserPrincipalFromMSALAuthToken(authToken)
	}

	profile, err := h.authService.GetProfile(userPrincipal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *AuthHandler) GetAllProfiles(c *gin.Context) {
	profiles, err := h.authService.GetAllProfiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profiles)
}

func (h *AuthHandler) AddRole(c *gin.Context) {
	userPrincipal := c.Param("userPrincipal")
	role := c.Param("role")
	err := h.authService.AddRole(userPrincipal, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *AuthHandler) AddDefaultProfile(c *gin.Context) {

	// Get the auth token from the request header
	authToken := c.GetHeader("Authorization")

	// Remove Bearer from the authToken
	authToken = strings.Split(authToken, "Bearer ")[1]

	userPrincipal, _ := helper.GetUserPrincipalFromMSALAuthToken(authToken)
	role := "user"

	slog.Info("Adding default role: " + role + " for user: " + userPrincipal)
	err := h.authService.AddRole(userPrincipal, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *AuthHandler) DeleteRole(c *gin.Context) {
	userPrincipal := c.Param("userPrincipal")
	role := c.Param("role")
	err := h.authService.DeleteRole(userPrincipal, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
