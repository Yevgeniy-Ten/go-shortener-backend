package lib

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateShortID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var shortID strings.Builder
	randSource := rand.NewSource(time.Now().UnixNano())
	r := rand.New(randSource)
	for i := 0; i < 8; i++ {
		shortID.WriteByte(charset[r.Intn(len(charset))])
	}

	return shortID.String()
}
