package middleware

import (
	"actlabs-auth/entity"
	"actlabs-auth/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		slog.Debug("Middleware: AuthRequired")

		// Get the auth token from the request header
		authToken := c.GetHeader("Authorization")

		// Remove Bearer from the authToken
		authToken = strings.Split(authToken, "Bearer ")[1]

		if authToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no auth token provided"})
			return
		}

		// Ensure that the token is issued by AAD.
		// isAADToken, err := helper.EnsureAADIssuer(authToken)
		// if err != nil || !isAADToken {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		// 	return
		// }

		c.Next()
	}
}

func AdminRequired(authService entity.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Debug("Middleware: AdminRequired")

		// Get the auth token from the request header
		authToken := c.GetHeader("Authorization")

		// Remove Bearer from the authToken
		authToken = strings.Split(authToken, "Bearer ")[1]

		// Ensure that the token is issued by AAD.
		// isAADToken, err := helper.EnsureAADIssuer(authToken)
		// if err != nil || !isAADToken {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		// 	return
		// }

		callingUserPrincipal, err := helper.GetUserPrincipalFromMSALAuthToken(authToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Allow all authenticated users to add 'user' role.
		role := c.Param("role")
		if role == "user" && c.Request.Method == "POST" {
			c.Next()
			return
		}

		// Get the roles for the calling user
		roles, err := authService.GetRoles(callingUserPrincipal)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Check if the calling user has the admin role
		if !helper.Contains(roles.Roles, "admin") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user is not an admin"})
			return
		}

		c.Next()
	}
}

func MentorRequired(authService entity.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Debug("Middleware: MentorRequired")

		// Get the auth token from the request header
		authToken := c.GetHeader("Authorization")

		// Remove Bearer from the authToken
		authToken = strings.Split(authToken, "Bearer ")[1]

		// Ensure that the token is issued by AAD.
		// isAADToken, err := helper.EnsureAADIssuer(authToken)
		// if err != nil || !isAADToken {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		// 	return
		// }

		callingUserPrincipal, err := helper.GetUserPrincipalFromMSALAuthToken(authToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Get the roles for the calling user
		roles, err := authService.GetRoles(callingUserPrincipal)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Check if the calling user has the mentor role
		if !helper.Contains(roles.Roles, "mentor") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user is not an mentor"})
			return
		}

		c.Next()
	}
}
