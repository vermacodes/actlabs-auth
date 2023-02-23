package helper

import (
	"crypto/rand"
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
