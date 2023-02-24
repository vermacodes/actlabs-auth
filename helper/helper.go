package helper

import (
	"crypto/rand"
	"strings"
	"unsafe"
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
