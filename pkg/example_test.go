package pkg_test

import (
	"fmt"
	"shorter/pkg"
)

func ExampleValidateURL() {
	testURLS := []string{
		"https://example.com/path?query=123",
		"http://example.com/path?query=123",
		"ftp://example.com/path?query=123",
		"example.com/path?query=123",
		"example.com",
		"example",
	}

	for _, url := range testURLS {
		if pkg.ValidateURL(url) {
			fmt.Println("URL is valid")
		} else {
			fmt.Println("URL is not valid")
		}
	}
	// Output:
	// URL is valid
	// URL is valid
	// URL is valid
	// URL is not valid
	// URL is not valid
	// URL is not valid
}
