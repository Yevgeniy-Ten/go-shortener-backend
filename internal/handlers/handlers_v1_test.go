package handlers_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"shorter/internal/domain"
	"shorter/internal/gzipper"
	"shorter/internal/handlers"
	"shorter/internal/urlstorage"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUrlHandlers(t *testing.T) {
	type value struct {
		value       string
		contentType string
	}
	type statusCodeCheck struct {
		statusCode int
	}
	createTests := []struct {
		name  string
		url   string
		value value
		want  statusCodeCheck
	}{
		{"Create Url", "/", value{"http://practicum.yandex.ru", "text/plain"}, statusCodeCheck{http.StatusCreated}},
		{"Create Url", "/", value{"http://practicum.yandex.ru", "text/plain"}, statusCodeCheck{http.StatusCreated}},
		{"Create Url", "/", value{"http://practicum.yandex.ru", "text/plain"}, statusCodeCheck{http.StatusCreated}},
		{
			"JSON CREATE", "/api/shorten", value{`{"url":"http://practicum.yandex.ru"}`,
				"application/json",
			},
			statusCodeCheck{http.StatusCreated},
		},
	}
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
	r := h.CreateRouter(gzipper.RequestResponseGzipMiddleware())

	createURL := func(value, contentType, url string) *http.Response {
		request := httptest.NewRequest("POST", url, strings.NewReader(value))
		request.Header.Set("Content-Type", contentType)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, request)
		result := recorder.Result()
		return result
	}

	for _, tt := range createTests {
		t.Run(tt.name, func(t *testing.T) {
			result := createURL(tt.value.value, tt.value.contentType, tt.url)
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			defer result.Body.Close()
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
			result := createURL(tt.value.value, tt.value.contentType, "/")
			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			request := httptest.NewRequest("GET", string(body), nil)
			err = result.Body.Close()
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, request)
			result = recorder.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))
			defer result.Body.Close()
		})
		t.Run("send_gzip", func(t *testing.T) {
			body := "http://practicum.yandex.ru"
			var buf bytes.Buffer
			gzipWriter := gzip.NewWriter(&buf)

			_, err := gzipWriter.Write([]byte(body))
			if err != nil {
				t.Fatalf("Ошибка при записи в gzipWriter: %v", err)
			}
			err = gzipWriter.Close()
			if err != nil {
				t.Fatalf("Ошибка при закрытии gzipWriter: %v", err)
			}
			request := httptest.NewRequest("POST", "/", &buf)
			request.Header.Set("Accept-Encoding", "gzip")
			request.Header.Set("Content-Type", "text/plain")
			request.Header.Set("Content-Encoding", "gzip")
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, request)
			result := recorder.Result()
			defer result.Body.Close()
			assert.Equal(t, http.StatusCreated, result.StatusCode)
		})
	}
}
