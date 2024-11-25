package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"shorter/internal/app/handlers"
	"strings"
	"testing"
)

func TestUrlHandler(t *testing.T) {
	type value struct {
		value       string
		contentType string
	}
	type statusCodeCheck struct {
		statusCode int
	}
	createTests := []struct { // добавляем слайс тестов
		name  string
		value value
		want  statusCodeCheck
	}{
		{"Check Content Type Url", value{"http://practicum.yandex.ru", "text/html"}, statusCodeCheck{400}},
		{"Create Url", value{"http://practicum.yandex.ru", "text/plain"}, statusCodeCheck{201}},
		{"Wrong Url", value{"htt://practicum.yandex.ru", "text/plain"}, statusCodeCheck{400}},
	}

	var createUrl func(string, string) *http.Response

	createUrl = func(value, contentType string) *http.Response {
		request := httptest.NewRequest("POST", "/", strings.NewReader(value))
		request.Header.Set("Content-Type", contentType)
		recorder := httptest.NewRecorder()
		handlers.URLHandler(recorder, request)
		result := recorder.Result()
		return result
	}

	for _, tt := range createTests { // перебираем все тесты
		t.Run(tt.name, func(t *testing.T) { // запускаем тест
			result := createUrl(tt.value.value, tt.value.contentType)
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}

	type redirectCheck struct {
		statusCodeCheck
		location string
	}
	redirectTests := []struct {
		name  string
		value value
		want  redirectCheck
	}{
		{"Redirect", value{"http://practicum.yandex.ru", "text/plain"}, redirectCheck{statusCodeCheck{307}, "http://practicum.yandex.ru"}},
		{"Redirect google", value{"https://google.com", "text/plain"}, redirectCheck{statusCodeCheck{307}, "https://google.com"}},
	}
	for _, tt := range redirectTests {
		t.Run(tt.name, func(t *testing.T) {
			result := createUrl(tt.value.value, tt.value.contentType)
			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			request := httptest.NewRequest("GET", string(body), nil)
			recorder := httptest.NewRecorder()
			handlers.URLHandler(recorder, request)
			result = recorder.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
		})
	}

}
