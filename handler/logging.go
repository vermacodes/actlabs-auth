package handler

import (
	"actlabs-auth/entity"
	"actlabs-auth/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type loggingHandler struct {
	loggingService entity.LoggingService
}

func NewLoggingHandler(r *gin.RouterGroup, loggingService entity.LoggingService) {
	handler := &loggingHandler{
		loggingService: loggingService,
	}
	r.POST("/logging/operation", handler.OperationRecord)
}

func (l *loggingHandler) OperationRecord(c *gin.Context) {
	// Get the auth token from the request header
	authToken := c.GetHeader("Authorization")

	// Remove Bearer from the authToken
	authToken = strings.Split(authToken, "Bearer ")[1]

	userPrincipal, _ := helper.GetUserPrincipalFromMSALAuthToken(authToken)

	var operation entity.Operation
	if err := c.BindJSON(&operation); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if err := l.loggingService.OperationRecord(operation, userPrincipal); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusCreated)
}
