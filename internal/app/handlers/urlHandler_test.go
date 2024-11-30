package handlers_test

import (
	"github.com/stretchr/testify/assert"
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
	type want struct {
		statusCode int
	}
	createTests := []struct { // добавляем слайс тестов
		name  string
		value value
		want  want
	}{
		{"Check Content Type Url", value{"http://practicum.yandex.ru", "text/html"}, want{400}},
		{"Create Url", value{"http://practicum.yandex.ru", "text/plain"}, want{201}},
		{"Wrong Url", value{"htt://practicum.yandex.ru", "text/plain"}, want{400}},
		{"Failed test XD ", value{"htt://pract.ru", "text/plain"}, want{201}},
	}
	for _, tt := range createTests { // перебираем все тесты
		t.Run(tt.name, func(t *testing.T) { // запускаем тест
			request := httptest.NewRequest("POST", "/", strings.NewReader(tt.value.value))
			request.Header.Set("Content-Type", tt.value.contentType)
			recorder := httptest.NewRecorder()
			handlers.URLHandler(recorder, request)
			result := recorder.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)

		})
	}

}
