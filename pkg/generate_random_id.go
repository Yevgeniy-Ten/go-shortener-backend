package pkg

import (
	"crypto/rand"
	"strings"
)

func GenerateShortID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8
	var shortID strings.Builder

	shortID.Grow(length)

	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("failed to generate random bytes")
	}

	for _, b := range randomBytes {
		shortID.WriteByte(charset[b%byte(len(charset))])
	}

	return shortID.String()
}
