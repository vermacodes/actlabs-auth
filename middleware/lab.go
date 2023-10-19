package middleware

import (
	"actlabs-auth/entity"
	"actlabs-auth/helper"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func UpdateCredits() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Debug("Middleware: UpdateCredits")

		// Get the auth token from the request header
		authToken := c.GetHeader("Authorization")

		// Remove Bearer from the authToken
		authToken = strings.Split(authToken, "Bearer ")[1]

		callingUserPrincipal, err := helper.GetUserPrincipalFromMSALAuthToken(authToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// get lab from the payload
		lab := entity.LabType{}
		if err := c.Bind(&lab); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// update credits
		if lab.Id == "" {
			lab.CreatedBy = callingUserPrincipal
			lab.CreatedOn = helper.GetTodaysDateString()
		} else {
			lab.UpdatedBy = callingUserPrincipal
			lab.UpdatedOn = helper.GetTodaysDateString()
		}

		slog.Debug("Updated lab", "lab", lab)

		// update request payload
		marshaledLab, err := json.Marshal(lab)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create a new request based on the existing request
		newRequest := c.Request.Clone(c.Request.Context())
		newRequest.Body = io.NopCloser(bytes.NewReader(marshaledLab))
		newRequest.ContentLength = int64(len(marshaledLab))

		// Replace the current request with the new request
		c.Request = newRequest

		slog.Debug("Updated request", "request", c.Request)

		c.Next()
	}
}
