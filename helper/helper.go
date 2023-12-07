package helper

import (
	"actlabs-auth/entity"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"unsafe"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
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

	userPrincipal, ok := tokenJSON["preferred_username"].(string)
	if !ok {
		err := errors.New("user principal name not found in token")
		slog.Error("user principal name not found in token", err)
		return "", err
	}

	return userPrincipal, nil
}

func VerifyToken(tokenString string) (bool, error) {

	// Drop the Bearer prefix if it exists
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.Split(tokenString, "Bearer ")[1]
	}

	keySet, err := jwk.Fetch(context.TODO(), "https://login.microsoftonline.com/common/discovery/v2.0/keys")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwa.RS256.String() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %v not found", kid)
		}

		publicKey := &rsa.PublicKey{}
		err = keys.Raw(publicKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key")
		}

		return publicKey, nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		err := errors.New("token is not valid")
		slog.Error("token is not valid", err)
		return false, err
	}

	// Get the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, errors.New("invalid claims")
	}

	// check the audience
	aud, ok := claims["aud"].(string)
	if !ok {
		return false, errors.New("invalid audience")
	}
	if aud != os.Getenv("AUTH_TOKEN_AUD") {
		return false, errors.New("invalid audience")
	}

	// Check the issuer
	iss, ok := claims["iss"].(string)
	if !ok {
		return false, errors.New("not able to get issuer from claims")
	}
	if iss != os.Getenv("AUTH_TOKEN_ISS") {
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

// Return today's date in the format yyyy-mm-dd as string
func GetTodaysDateString() string {
	return time.Now().Format("2006-01-02")
}

// Return today's date and time in the format yyyy-mm-dd hh:mm:ss as string
func GetTodaysDateTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// Return today's date and time in ISO format as string
func GetTodaysDateTimeISOString() string {
	return time.Now().Format(time.RFC3339)
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
