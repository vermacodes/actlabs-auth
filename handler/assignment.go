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

	r.GET("/assignment/labs", handler.GetAllLabsRedacted)
	r.GET("/assignment/labs/my", handler.GetMyAssignedLabsRedacted)
	r.GET("/assignment/my", handler.GetMyAssignments)
	r.POST("/assignment/my", handler.CreateMyAssignments)
	r.DELETE("/assignment/my", handler.DeleteMyAssignments)
}

func NewAssignmentHandlerMentorRequired(r *gin.RouterGroup, service entity.AssignmentService) {
	handler := &assignmentHandler{
		assignmentService: service,
	}

	r.GET("/assignment", handler.GetAllAssignments)
	r.GET("/assignment/lab/:labId", handler.GetAssignmentsByLabId)
	r.GET("/assignment/user/:userId", handler.GetAssignmentsByUserId)
	r.POST("/assignment", handler.CreateAssignments)
	r.DELETE("/assignment", handler.DeleteAssignments)
}

func (a *assignmentHandler) GetAllAssignments(c *gin.Context) {
	assignments, err := a.assignmentService.GetAllAssignments()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, assignments)
}

func (a *assignmentHandler) GetAllLabsRedacted(c *gin.Context) {
	labs, err := a.assignmentService.GetAllLabsRedacted()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusOK, labs)
}

func (a *assignmentHandler) GetMyAssignedLabsRedacted(c *gin.Context) {

	// Get the auth token from the request header
	authToken := c.GetHeader("Authorization")

	// Remove Bearer from the authToken
	authToken = strings.Split(authToken, "Bearer ")[1]
	//Get the user principal from the auth token
	userId, _ := helper.GetUserPrincipalFromMSALAuthToken(authToken)

	labs, err := a.assignmentService.GetAssignedLabsRedactedByUserId(userId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, labs)
}

func (a *assignmentHandler) GetAssignedLabsRedactedByUserId(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "my" {
		// Get the auth token from the request header
		authToken := c.GetHeader("Authorization")

		// Remove Bearer from the authToken
		authToken = strings.Split(authToken, "Bearer ")[1]
		//Get the user principal from the auth token
		userId, _ = helper.GetUserPrincipalFromMSALAuthToken(authToken)
	}

	labs, err := a.assignmentService.GetAssignedLabsRedactedByUserId(userId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, labs)
}

func (a *assignmentHandler) GetAssignmentsByLabId(c *gin.Context) {
	labId := c.Param("labId")
	assignments, err := a.assignmentService.GetAssignmentsByLabId(labId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, assignments)
}

func (a *assignmentHandler) GetAssignmentsByUserId(c *gin.Context) {
	userId := c.Param("userId")
	assignments, err := a.assignmentService.GetAssignmentsByUserId(userId)
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

	assignments, err := a.assignmentService.GetAssignmentsByUserId(userPrincipal)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, assignments)
}

func (a *assignmentHandler) CreateMyAssignments(c *gin.Context) {
	bulkAssignment := entity.BulkAssignment{}
	if err := c.Bind(&bulkAssignment); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Get the auth token from the request header
	authToken := c.GetHeader("Authorization")

	// Remove Bearer from the authToken
	authToken = strings.Split(authToken, "Bearer ")[1]
	//Get the user principal from the auth token
	userPrincipal, _ := helper.GetUserPrincipalFromMSALAuthToken(authToken)

	// Sanitizing to make sure that the user is not creating assignments for other users.
	for _, userId := range bulkAssignment.UserIds {
		if userId != userPrincipal {
			c.Status(http.StatusBadRequest)
			return
		}
	}

	if err := a.assignmentService.CreateAssignments(bulkAssignment.UserIds, bulkAssignment.LabIds, userPrincipal); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (a *assignmentHandler) CreateAssignments(c *gin.Context) {
	bulkAssignment := entity.BulkAssignment{}
	if err := c.Bind(&bulkAssignment); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Get the auth token from the request header
	authToken := c.GetHeader("Authorization")

	// Remove Bearer from the authToken
	authToken = strings.Split(authToken, "Bearer ")[1]
	//Get the user principal from the auth token
	userPrincipal, _ := helper.GetUserPrincipalFromMSALAuthToken(authToken)

	if err := a.assignmentService.CreateAssignments(bulkAssignment.UserIds, bulkAssignment.LabIds, userPrincipal); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (a *assignmentHandler) DeleteMyAssignments(c *gin.Context) {
	assignments := []string{}
	if err := c.Bind(&assignments); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Get the auth token from the request header
	authToken := c.GetHeader("Authorization")

	// Remove Bearer from the authToken
	authToken = strings.Split(authToken, "Bearer ")[1]
	//Get the user principal from the auth token
	userPrincipal, _ := helper.GetUserPrincipalFromMSALAuthToken(authToken)

	// Sanitizing to make sure that the user is not deleting assignments for other users.
	for _, assignment := range assignments {
		if !strings.HasPrefix(assignment, userPrincipal) {
			c.Status(http.StatusBadRequest)
			return
		}
	}

	if err := a.assignmentService.DeleteAssignments(assignments); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *assignmentHandler) DeleteAssignments(c *gin.Context) {
	assignments := []string{}
	if err := c.Bind(&assignments); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := a.assignmentService.DeleteAssignments(assignments); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
