package pkg

import "testing"

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
