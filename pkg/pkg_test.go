package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkGenerateShortID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenerateShortID()
	}
}

func BenchmarkValidateURL(b *testing.B) {
	testURL := "https://example.com/path?query=123"

	for i := 0; i < b.N; i++ {
		_ = ValidateURL(testURL)
	}
}

func TestGenerateShortID(t *testing.T) {
	for i := 0; i < 20; i++ {
		got := GenerateShortID()
		assert.Equal(t, 8, len(got))
	}
}
