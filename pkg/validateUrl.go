package pkg

func ValidateURL(url string) bool {
	return url[:4] == "http"
}
