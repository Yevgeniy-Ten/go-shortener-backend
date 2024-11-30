package main

import (
	"log"
	"net/http"
	handlers "shorter/internal/app/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.URLHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
