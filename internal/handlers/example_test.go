package handlers_test

import (
	"fmt"
	"net/http/httptest"
	"shorter/internal/domain"
	"shorter/internal/handlers"
	"shorter/internal/urlstorage"
	"strings"
)

func Example() {
	s := urlstorage.New(nil)
	h := handlers.NewHandler(
		&handlers.Config{
			ServerAddr: "http://localhost:8080",
		},
		domain.Storage{
			URLS: s,
		},
		nil,
	)
	r := h.CreateRouter()
	request := httptest.NewRequest("POST", "/", strings.NewReader("http://practicum.yandex.ru"))
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)
	result := recorder.Result()
	defer result.Body.Close()
	fmt.Println(result.StatusCode)
	shortURL := recorder.Body.String()
	request = httptest.NewRequest("GET", shortURL, nil)
	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, request)
	result = recorder.Result()
	defer result.Body.Close()
	fmt.Println(result.StatusCode, result.Header.Get("Location"))
	// Output: 201
	// 307 http://practicum.yandex.ru
}
