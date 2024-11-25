package main

import (
	"fmt"
	"net/http"
)

func main() {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println(req.URL)
			return nil
		},
	}
	response, _ := client.Get("http://ya.ru")
	for k, v := range response.Header {
		fmt.Println(k, v)
	}
	defer response.Body.Close()
}
