package middleware

import (
	"actlabs-auth/entity"
	"actlabs-auth/helper"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/slog"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		slog.Info("Middleware: AuthRequired")

		// Get the auth token from the request header
		authToken := c.GetHeader("Authorization")

		// Remove Bearer from the authToken
		authToken = strings.Split(authToken, "Bearer ")[1]

		if authToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no auth token provided"})
			return
		}

		// Ensure that the token is issued by AAD.
		isAADToken, err := ensureAADIssuer(authToken)
		if err != nil || !isAADToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Next()
	}
}

func AdminRequired(authService entity.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("Middleware: AdminRequired")

		// Get the auth token from the request header
		authToken := c.GetHeader("Authorization")

		// Remove Bearer from the authToken
		authToken = strings.Split(authToken, "Bearer ")[1]

		// Ensure that the token is issued by AAD.
		isAADToken, err := ensureAADIssuer(authToken)
		if err != nil || !isAADToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		callingUserPrincipal, err := getUserPrincipalFromMSALAuthToken(authToken)
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
		slog.Info("Middleware: MentorRequired")

		// Get the auth token from the request header
		authToken := c.GetHeader("Authorization")

		// Remove Bearer from the authToken
		authToken = strings.Split(authToken, "Bearer ")[1]

		// Ensure that the token is issued by AAD.
		isAADToken, err := ensureAADIssuer(authToken)
		if err != nil || !isAADToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		callingUserPrincipal, err := getUserPrincipalFromMSALAuthToken(authToken)
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

// Helper Functions.

// Get user principal name from MSAL auth token.
func getUserPrincipalFromMSALAuthToken(token string) (string, error) {

	// Split the token into its parts
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) < 2 {
		println("invalid token format")
		return "", errors.New("invalid token format")
	}

	// Decode the token
	decodedToken, err := base64.StdEncoding.DecodeString(tokenParts[1] + strings.Repeat("=", (4-len(tokenParts[1])%4)%4))
	if err != nil {
		println("not able to decode token -> ", err.Error())
		return "", err
	}

	// Extract the user principal name from the decoded token
	var tokenJSON map[string]interface{}
	err = json.Unmarshal(decodedToken, &tokenJSON)
	if err != nil {
		println("not able to unmarshal token -> ", err.Error())
		return "", err
	}

	userPrincipal, ok := tokenJSON["upn"].(string)
	if !ok {
		println("user principal name not found in token")
		return "", errors.New("user principal name not found in token")
	}

	return userPrincipal, nil
}

// Ensure that the token is issued by AAD.
func ensureAADIssuer(tokenString string) (bool, error) {

	publicKeyString, err := getPublicKey(tokenString)
	if err != nil {
		return false, err
	}

	publicKey := "-----BEGIN CERTIFICATE-----\n " + publicKeyString + "\n-----END CERTIFICATE-----"

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Decode public key
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %v", err)
		}

		return pubKey, nil

	})

	if err != nil {
		println("not able to parse token -> ", err.Error())
		return false, err
	}

	// Check if the token is valid
	if !token.Valid {
		return false, errors.New("invalid token")
	}

	// Get the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, errors.New("invalid claims")
	}

	// Check the issuer
	iss, ok := claims["iss"].(string)
	if !ok {
		return false, errors.New("invalid issuer")
	}
	if iss != "https://sts.windows.net/72f988bf-86f1-41af-91ab-2d7cd011db47/" {
		return false, errors.New("invalid issuer")
	}

	// Check the expiration time
	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, errors.New("invalid expiration time")
	}
	if time.Now().Unix() > int64(exp) {
		return false, errors.New("token has expired")
	}

	return true, nil
}

// Get the public key from the well-known endpoint.
func getKid(token string) (string, error) {

	// Split the token
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errors.New("invalid token")
	}

	// Decode the header
	header, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", err
	}

	// Get the kid from the header
	var headerJSON map[string]interface{}
	err = json.Unmarshal(header, &headerJSON)
	if err != nil {
		return "", err
	}
	kid, ok := headerJSON["kid"].(string)
	if !ok {
		return "", errors.New("kid not found in token")
	}

	return kid, nil
}

// Get the public key from the well-known endpoint.
func getPublicKeyFromWellKnown(kid string) (string, error) {

	// Get the well-known endpoint
	resp, err := http.Get("https://login.microsoftonline.com/common/discovery/keys")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode the response
	var wellKnownJSON map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&wellKnownJSON)
	if err != nil {
		return "", err
	}

	// Get the keys from the response
	keys, ok := wellKnownJSON["keys"].([]interface{})
	if !ok {
		return "", errors.New("keys not found in well-known response")
	}

	// Find the key with the matching kid
	for _, key := range keys {
		keyMap, ok := key.(map[string]interface{})
		if !ok {
			return "", errors.New("invalid key")
		}
		if keyMap["kid"] == kid {
			x5c, ok := keyMap["x5c"].([]interface{})
			if !ok {
				return "", errors.New("invalid x5c value")
			}
			var x5cStrings []string
			for _, v := range x5c {
				x5cStrings = append(x5cStrings, v.(string))
			}
			return strings.Join(x5cStrings, ""), nil
		}
	}

	return "", errors.New("key not found")
}

// Get the public key from the well-known endpoint.
func getPublicKey(token string) (string, error) {

	// Get the kid from the token
	kid, err := getKid(token)
	if err != nil {
		return "", err
	}

	// Get the public key from the well-known endpoint
	publicKey, err := getPublicKeyFromWellKnown(kid)
	if err != nil {
		return "", err
	}

	return publicKey, nil
}
