package helper

import (
	"actlabs-auth/entity"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
	"unsafe"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/slog"
)

var alphabet = []byte("abcdefghijklmnopqrstuvwxyz0123456789")

func Generate(length int) string {
	// Generate a alphanumeric string of length length.

	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = alphabet[b[i]%byte(len(alphabet))]
	}
	return *(*string)(unsafe.Pointer(&b))
}

// Function to convert a slice of strings to a single string delimited by a comma
func SliceToString(s []string) string {
	return strings.Join(s, ",")
}

// Function to convert a string delimited by a comma to a slice of strings
func StringToSlice(s string) []string {
	return strings.Split(s, ",")
}

// Function to check if a string is in a slice of strings
func Contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// SlicesAreEqual checks if two slices are equal.
func SlicesAreEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func GetUserPrincipalFromMSALAuthToken(token string) (string, error) {

	// Split the token into its parts
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) < 2 {
		err := errors.New("invalid token format")
		slog.Error("invalid token format", err)
		return "", err
	}

	// Decode the token
	decodedToken, err := base64.StdEncoding.DecodeString(tokenParts[1] + strings.Repeat("=", (4-len(tokenParts[1])%4)%4))
	if err != nil {
		slog.Error("not able to decode token -> ", err)
		return "", err
	}

	// Extract the user principal name from the decoded token
	var tokenJSON map[string]interface{}
	err = json.Unmarshal(decodedToken, &tokenJSON)
	if err != nil {
		slog.Error("not able to unmarshal token -> ", err)
		return "", err
	}

	userPrincipal, ok := tokenJSON["upn"].(string)
	if !ok {
		err := errors.New("user principal name not found in token")
		slog.Error("user principal name not found in token", err)
		return "", err
	}

	return userPrincipal, nil
}

// Ensure that the token is issued by AAD.
func EnsureAADIssuer(tokenString string) (bool, error) {

	publicKeyString, err := GetPublicKey(tokenString)
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
		slog.Error("not able to parse token -> ", err)
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
func GetKid(token string) (string, error) {

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
func GetPublicKeyFromWellKnown(kid string) (string, error) {

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
func GetPublicKey(token string) (string, error) {

	// Get the kid from the token
	kid, err := GetKid(token)
	if err != nil {
		return "", err
	}

	// Get the public key from the well-known endpoint
	publicKey, err := GetPublicKeyFromWellKnown(kid)
	if err != nil {
		return "", err
	}

	return publicKey, nil
}

// Return today's date in the format yyyy-mm-dd as string
func GetTodaysDateString() string {
	return time.Now().Format("2006-01-02")
}

// Return today's date and time in the format yyyy-mm-dd hh:mm:ss as string
func GetTodaysDateTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// ConvertProfileToRecord converts a Profile to a ProfileRecord.
func ConvertProfileToRecord(profile entity.Profile) entity.ProfileRecord {
	return entity.ProfileRecord{
		PartitionKey:  "actlabs",             // this is a static value.
		RowKey:        profile.UserPrincipal, // UserPrincipal is the unique identifier for the user.
		ObjectId:      profile.ObjectId,
		UserPrincipal: profile.UserPrincipal,
		DisplayName:   profile.DisplayName,
		ProfilePhoto:  profile.ProfilePhoto,
		Roles:         strings.Join(profile.Roles, ","),
	}
}

// ConvertRecordToProfile converts a ProfileRecord to a Profile.
func ConvertRecordToProfile(record entity.ProfileRecord) entity.Profile {
	return entity.Profile{
		ObjectId:      record.ObjectId,
		UserPrincipal: record.UserPrincipal,
		DisplayName:   record.DisplayName,
		ProfilePhoto:  record.ProfilePhoto,
		Roles:         strings.Split(record.Roles, ","),
	}
}
