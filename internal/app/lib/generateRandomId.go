package generateRandomId

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateShortId() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var shortId strings.Builder
	randSource := rand.NewSource(time.Now().UnixNano())
	r := rand.New(randSource)
	for i := 0; i < 8; i++ {
		shortId.WriteByte(charset[r.Intn(len(charset))])
	}

	return shortId.String()
}
