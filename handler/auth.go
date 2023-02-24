package handler

import (
	"actlabs-auth/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService entity.AuthService
}

func NewAuthHandler(r *gin.RouterGroup, authService entity.AuthService) {
	handler := &AuthHandler{
		authService: authService,
	}

	r.GET("/roles/:userPrincipal", handler.GetRoles)
}

func NewAdminAuthHandler(r *gin.RouterGroup, authService entity.AuthService) {
	handler := &AuthHandler{
		authService: authService,
	}

	r.POST("/roles/:userPrincipal/:role", handler.AddRole)
	r.DELETE("/roles/:userPrincipal/:role", handler.DeleteRole)
}

func (h *AuthHandler) GetRoles(c *gin.Context) {
	userPrincipal := c.Param("userPrincipal")
	roles, err := h.authService.GetRoles(userPrincipal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
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
