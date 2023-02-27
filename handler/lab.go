package handler

import (
	"net/http"

	"actlabs-auth/entity"

	"github.com/gin-gonic/gin"
)

type labHandler struct {
	labService entity.LabService
}

func NewLabHandler(r *gin.RouterGroup, labService entity.LabService) {
	handler := &labHandler{
		labService: labService,
	}
	r.GET("/lab/public/:typeOfLab", handler.GetPublicLabs)
}

func NewLabHandlerMentorRequired(r *gin.RouterGroup, labService entity.LabService) {
	handler := &labHandler{
		labService: labService,
	}
	r.POST("/lab/protected", handler.AddLab)
	r.DELETE("/lab/protected", handler.DeleteLab)
	r.GET("/lab/protected/:typeOfLab", handler.GetProtectedLabs)
}

func (l *labHandler) GetPublicLabs(c *gin.Context) {
	typeOfLab := c.Param("typeOfLab")

	// These labs are protected, use protected API
	if typeOfLab == "mockcases" || typeOfLab == "labexercises" {
		c.Status(http.StatusBadRequest)
		return
	}

	labs, err := l.labService.GetPublicLabs(typeOfLab)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, labs)
}

func (l *labHandler) GetProtectedLabs(c *gin.Context) {
	typeOfLab := c.Param("typeOfLab")
	labs, err := l.labService.GetPublicLabs(typeOfLab)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, labs)
}

func (l *labHandler) AddLab(c *gin.Context) {
	lab := entity.LabType{}
	if err := c.Bind(&lab); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := l.labService.AddPublicLab(lab); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (l *labHandler) DeleteLab(c *gin.Context) {
	lab := entity.LabType{}
	if err := c.Bind(&lab); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := l.labService.DeletePublicLab(lab); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
