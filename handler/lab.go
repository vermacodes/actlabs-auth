package handler

import (
	"net/http"
	"strings"

	"actlabs-auth/entity"
	"actlabs-auth/helper"

	"github.com/gin-gonic/gin"
)

type labHandler struct {
	labService entity.LabService
}

// Authenticated user.
func NewLabHandler(r *gin.RouterGroup, labService entity.LabService) {
	handler := &labHandler{
		labService: labService,
	}
	// all private lab operations.
	r.GET("/lab/private/:typeOfLab", handler.GetLabs)
	r.POST("/lab/private", handler.UpsertLab)
	r.DELETE("/lab/private/:typeOfLab/:labId", handler.DeleteLab)
	r.GET("/lab/private/versions/:typeOfLab/:labId", handler.GetLabVersions)

	// public lab read-only operations.
	r.GET("/lab/public/:typeOfLab", handler.GetLabs)
	r.GET("/lab/public/versions/:typeOfLab/:labId", handler.GetLabVersions)
}

// Authenticated user with 'contributor' role.
func NewLabHandlerContributorRequired(r *gin.RouterGroup, labService entity.LabService) {
	handler := &labHandler{
		labService: labService,
	}

	// public lab mutable operations.
	r.POST("/lab/public", handler.UpsertLab)
	r.DELETE("/lab/public/:typeOfLab/:labId", handler.DeleteLab)
}

// Authenticated user with 'mentor' role.
func NewLabHandlerMentorRequired(r *gin.RouterGroup, labService entity.LabService) {
	handler := &labHandler{
		labService: labService,
	}

	// all protected lab operations.
	r.POST("/lab/protected", handler.UpsertLab)
	r.GET("/lab/protected/:typeOfLab", handler.GetLabs)
	r.GET("/lab/protected/versions/:typeOfLab/:labId", handler.GetLabVersions)
	r.DELETE("/lab/protected/:typeOfLab/:labId", handler.DeleteLab)
}

func (l *labHandler) GetLabs(c *gin.Context) {
	typeOfLab := c.Param("typeOfLab")

	var labs []entity.LabType
	var err error

	switch {
	case validateLabType(typeOfLab, entity.PrivateLab):
		// Get the auth token from the request header
		authToken := c.GetHeader("Authorization")
		// Remove Bearer from the authToken
		authToken = strings.Split(authToken, "Bearer ")[1]
		userId, _ := helper.GetUserPrincipalFromMSALAuthToken(authToken)
		labs, err = l.labService.GetPrivateLabs(typeOfLab, userId)
	case validateLabType(typeOfLab, entity.PublicLab):
		labs, err = l.labService.GetPublicLab(typeOfLab)
	case validateLabType(typeOfLab, entity.ProtectedLabs):
		labs, err = l.labService.GetProtectedLabs(typeOfLab)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lab type: " + typeOfLab})
		return
	}

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, labs)
}

func (l *labHandler) UpsertLab(c *gin.Context) {
	lab := entity.LabType{}
	if err := c.Bind(&lab); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var err error

	switch {
	case validateLabType(lab.Type, entity.PrivateLab):
		err = l.labService.UpsertPrivateLab(lab)
	case validateLabType(lab.Type, entity.PublicLab):
		err = l.labService.UpsertPublicLab(lab)
	case validateLabType(lab.Type, entity.ProtectedLabs):
		err = l.labService.UpsertProtectedLab(lab)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lab type: " + lab.Type})
		return
	}

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (l *labHandler) DeleteLab(c *gin.Context) {
	typeOfLab := c.Param("typeOfLab")
	labId := c.Param("labId")

	var err error

	// Get the auth token from the request header
	authToken := c.GetHeader("Authorization")
	// Remove Bearer from the authToken
	authToken = strings.Split(authToken, "Bearer ")[1]
	userId, _ := helper.GetUserPrincipalFromMSALAuthToken(authToken)

	switch {
	case validateLabType(typeOfLab, entity.PrivateLab):
		err = l.labService.DeletePrivateLab(typeOfLab, labId, userId)
	case validateLabType(typeOfLab, entity.PublicLab):
		err = l.labService.DeletePublicLab(typeOfLab, labId, userId)
	case validateLabType(typeOfLab, entity.ProtectedLabs):
		err = l.labService.DeleteProtectedLab(typeOfLab, labId)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lab type: " + typeOfLab})
		return
	}

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (l *labHandler) GetLabVersions(c *gin.Context) {
	typeOfLab := c.Param("typeOfLab")
	labId := c.Param("labId")

	var labs []entity.LabType
	var err error

	switch {
	case validateLabType(typeOfLab, entity.PrivateLab):
		// Get the auth token from the request header
		authToken := c.GetHeader("Authorization")
		// Remove Bearer from the authToken
		authToken = strings.Split(authToken, "Bearer ")[1]
		userId, _ := helper.GetUserPrincipalFromMSALAuthToken(authToken)
		labs, err = l.labService.GetPrivateLabVersions(typeOfLab, labId, userId)
	case validateLabType(typeOfLab, entity.PublicLab):
		labs, err = l.labService.GetPublicLabVersions(typeOfLab, labId)
	case validateLabType(typeOfLab, entity.ProtectedLabs):
		labs, err = l.labService.GetProtectedLabVersions(typeOfLab, labId)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lab type: " + typeOfLab})
		return
	}

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, labs)
}

func validateLabType(typeOfLab string, validTypes []string) bool {
	for _, t := range validTypes {
		if typeOfLab == t {
			return true
		}
	}
	return false
}
