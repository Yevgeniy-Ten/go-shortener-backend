package pkg

import (
	"net/url"
)

// ValidateURL validates the URL
func ValidateURL(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}
	return true
}
