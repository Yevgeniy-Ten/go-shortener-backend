package lib

import "fmt"

func ValidateUrl(url string) bool {
	fmt.Println(url[:4])
	if url[:4] != "http" {
		return false
	}
	return true
}
