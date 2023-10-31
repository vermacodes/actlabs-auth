package handler

import (
	"net/http"
	"strings"

	"actlabs-auth/entity"
	"actlabs-auth/helper"

	"github.com/gin-gonic/gin"
)

type assignmentHandler struct {
	assignmentService entity.AssignmentService
}

func NewAssignmentHandler(r *gin.RouterGroup, service entity.AssignmentService) {
	handler := &assignmentHandler{
		assignmentService: service,
	}

	r.GET("/assignment", handler.GetAssignments)
	r.GET("/assignment/my", handler.GetMyAssignments)
	r.POST("/assignment", handler.CreateAssignment)
	r.DELETE("/assignment", handler.DeleteAssignment)
}

func (a *assignmentHandler) GetAssignments(c *gin.Context) {
	assignments, err := a.assignmentService.GetAssignments()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, assignments)
}

func (a *assignmentHandler) GetMyAssignments(c *gin.Context) {
	// Get the auth token from the request header
	authToken := c.GetHeader("Authorization")

	// Remove Bearer from the authToken
	authToken = strings.Split(authToken, "Bearer ")[1]
	//Get the user principal from the auth token
	userPrincipal, _ := helper.GetUserPrincipalFromMSALAuthToken(authToken)

	assignments, err := a.assignmentService.GetMyAssignments(userPrincipal)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, assignments)
}

func (a *assignmentHandler) CreateAssignment(c *gin.Context) {
	assignment := entity.Assignment{}
	if err := c.Bind(&assignment); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := a.assignmentService.CreateAssignment(assignment); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (a *assignmentHandler) DeleteAssignment(c *gin.Context) {
	assignment := entity.Assignment{}
	if err := c.Bind(&assignment); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := a.assignmentService.DeleteAssignment(assignment); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}
