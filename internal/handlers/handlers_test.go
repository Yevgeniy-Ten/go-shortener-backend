package handlers_test

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"shorter/internal/cookies"
	"shorter/internal/domain"
	"shorter/internal/handlers"
	"shorter/internal/handlers/mocks"
	"shorter/internal/logger"
	"strings"
	"testing"
)

func getInitial(t *testing.T) (*gomock.Controller, *gin.Engine, *handlers_mocks.MockuserStorage, *handlers_mocks.MockurlStorage, *handlers.Config, logger.ZapLogger, domain.Storage) {
	ctrl := gomock.NewController(t)
	mockStorage := handlers_mocks.NewMockuserStorage(ctrl)
	mockURLStorage := handlers_mocks.NewMockurlStorage(ctrl)
	cfg := &handlers.Config{
		ServerAddr:    "http://localhost:8080",
		DatabaseURL:   "soso",
		DatabaseError: false,
	}
	var l = logger.ZapLogger{
		Log: zap.NewNop(),
	}
	s := domain.Storage{
		User: mockStorage,
		URLS: mockURLStorage,
	}
	h := handlers.InitHandlers(cfg, s, &l)
	return ctrl, h, mockStorage, mockURLStorage, cfg, l, s

}
func TestInitHandlers(t *testing.T) {
	ctrl, _, _, _, cfg, l, s := getInitial(t)

	defer ctrl.Finish()

	t.Run("Handlers inited", func(t *testing.T) {

		h := handlers.InitHandlers(cfg, s, &l)
		assert.NotNil(t, h)
	})
}
func TestCreateRouter(t *testing.T) {
	ctrl, _, _, _, cfg, l, s := getInitial(t)
	defer ctrl.Finish()
	t.Run("Create routes created", func(t *testing.T) {
		h := handlers.NewHandler(
			cfg, s, &l)
		r := h.CreateRouter()
		assert.NotNil(t, r)
	})
}

func TestHandler_PostHandler(t *testing.T) {
	ctrl, h, mockStorage, mockURLStorage, _, _, _ := getInitial(t)
	defer ctrl.Finish()
	t.Run("PostHandler", func(t *testing.T) {
		mockStorage.EXPECT().Create().Return(1, nil)
		mockURLStorage.EXPECT().Save("http://practicum.yandex.ru", 1).Return("123", nil)
		request := httptest.NewRequest("POST", "/", strings.NewReader("http://practicum.yandex.ru"))
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		defer result.Body.Close()
		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.Equal(t, "http://localhost:8080/123", recorder.Body.String())
	})
	t.Run("PostHandler with error", func(t *testing.T) {
		mockStorage.EXPECT().Create().Return(0, nil)
		request := httptest.NewRequest("POST", "/", strings.NewReader("practicum.yandex.ru"))
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	})
	t.Run("PostHandler with error", func(t *testing.T) {
		mockStorage.EXPECT().Create().Return(1, nil)
		mockURLStorage.EXPECT().Save("http://practicum.yandex.ru", 1).Return("", errors.New("error"))
		request := httptest.NewRequest("POST", "/", strings.NewReader("http://practicum.yandex.ru"))
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	})

}

func TestHandler_GetHandler(t *testing.T) {

	ctrl, h, _, mockURLStorage, _, _, _ := getInitial(t)
	defer ctrl.Finish()
	t.Run("GetHandler", func(t *testing.T) {
		mockURLStorage.EXPECT().GetURL("123").Return("http://practicum.yandex.ru", nil)
		request := httptest.NewRequest("GET", "/123", nil)
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode)
		assert.Equal(t, "http://practicum.yandex.ru", result.Header.Get("Location"))
	})
	t.Run("GetHandler with error", func(t *testing.T) {
		mockURLStorage.EXPECT().GetURL("123").Return("", errors.New("error"))
		request := httptest.NewRequest("GET", "/123", nil)
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusNotFound, result.StatusCode)
	})
}
func TestHandler_ShortenURLHandler(t *testing.T) {
	ctrl, h, mockStorage, mockURLStorage, _, _, _ := getInitial(t)
	defer ctrl.Finish()
	t.Run("ShortenURLHandler", func(t *testing.T) {

		mockStorage.EXPECT().Create().Return(1, nil)
		mockURLStorage.EXPECT().Save("http://practicum.yandex.ru", 1).Return("123", nil)
		request := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(`{"url":"http://practicum.yandex.ru"}`))
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.Equal(t,
			"{\"result\":\"http://localhost:8080/123\"}", recorder.Body.String())
	})
	t.Run("ShortenURLHandler with error", func(t *testing.T) {
		mockStorage.EXPECT().Create().Return(0, nil)
		request := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(`{"url":"practicum.yandex.ru"}`))
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	})
	t.Run("ShortenURLHandler with error", func(t *testing.T) {
		mockStorage.EXPECT().Create().Return(1, nil)
		mockURLStorage.EXPECT().Save("http://practicum.yandex.ru", 1).Return("", errors.New("error"))
		request := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(`{"url":"http://practicum.yandex.ru"}`))
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	})
}

func TestHandler_GetUserUrls(t *testing.T) {
	ctrl, h, mockStorage, mockURLStorage, _, _, _ := getInitial(t)
	defer ctrl.Finish()
	t.Run("GetUserUrls", func(t *testing.T) {
		mockStorage.EXPECT().Create().Return(1, nil)
		mockURLStorage.EXPECT().GetUserURLs(
			1, "http://localhost:8080/").Return(
			[]domain.UserURLs{
				{OriginalURL: "http://localhost:8080/123", ShortURL: "123"},
			}, nil)
		request := httptest.NewRequest("GET", "/api/user/urls", nil)
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.JSONEq(t, `[{"original_url":"http://localhost:8080/123","short_url":"123"}]`, recorder.Body.String())
	})
	t.Run("GetUserUrls with error", func(t *testing.T) {
		mockStorage.EXPECT().Create().Return(0, nil)
		mockURLStorage.EXPECT().GetUserURLs(
			0, "http://localhost:8080/").Return(
			[]domain.UserURLs{}, errors.New("error"))
		request := httptest.NewRequest("GET", "/api/user/urls", nil)
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	})
	t.Run("GetUserUrls no urls", func(t *testing.T) {
		mockStorage.EXPECT().Create().Return(1, nil)
		mockURLStorage.EXPECT().GetUserURLs(
			1, "http://localhost:8080/").Return(
			[]domain.UserURLs{}, nil)
		request := httptest.NewRequest("GET", "/api/user/urls", nil)
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusNoContent, result.StatusCode)
	})
}
func TestHandler_DeleteMyUrls(t *testing.T) {
	ctrl, h, _, mockURLStorage, _, _, _ := getInitial(t)
	defer ctrl.Finish()
	t.Run("DeleteMyUrls", func(t *testing.T) {
		ids := []string{"1", "2", "3"}
		mockURLStorage.EXPECT().DeleteURLs(ids, 1).Return(nil).AnyTimes()
		request := httptest.NewRequest("DELETE", "/api/user/urls",
			strings.NewReader(`["1","2","3"]`))
		c, err := cookies.CreateCookie(1)
		assert.NoError(t, err)
		cookie := &http.Cookie{
			Name:  cookies.CookieName,
			Value: c,
			Path:  "/",
		}
		request.AddCookie(cookie)
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusAccepted, result.StatusCode)
	})
	t.Run("DeleteMyUrls with 401", func(t *testing.T) {
		request := httptest.NewRequest("DELETE", "/api/user/urls", nil)
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
	})
	t.Run("DeleteMyUrls with error", func(t *testing.T) {
		request := httptest.NewRequest("DELETE", "/api/user/urls",
			nil)
		c, err := cookies.CreateCookie(1)
		assert.NoError(t, err)
		cookie := &http.Cookie{
			Name:  cookies.CookieName,
			Value: c,
			Path:  "/",
		}
		request.AddCookie(cookie)
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	})
}

func TestHandler_ShortenURLSHandler(t *testing.T) {
	ctrl, h, _, mockURLStorage, _, _, _ := getInitial(t)
	defer ctrl.Finish()
	t.Run("ShortenURLSHandler", func(t *testing.T) {
		mockURLStorage.EXPECT().SaveBatch(gomock.Any(), 1).Return(nil)
		payload := `[{"correlation_id":"123","original_url":"http://practicum.yandex.ru"}]`
		request := httptest.NewRequest("POST", "/api/shorten/batch", strings.NewReader(payload))
		c, err := cookies.CreateCookie(1)
		assert.NoError(t, err)
		cookie := &http.Cookie{
			Name:  cookies.CookieName,
			Value: c,
			Path:  "/",
		}
		request.AddCookie(cookie)
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.JSONEq(t, `[{"correlation_id":"123","short_url":"http://localhost:8080/123"}]`, recorder.Body.String())
	})
	t.Run("ShortenURLSHandler with error", func(t *testing.T) {
		mockURLStorage.EXPECT().SaveBatch(gomock.Any(), 1).Return(errors.New("error"))
		payload := `[{"correlation_id":"123","original_url":"http://practicum.yandex.ru"}]`
		request := httptest.NewRequest("POST", "/api/shorten/batch", strings.NewReader(payload))
		c, err := cookies.CreateCookie(1)
		assert.NoError(t, err)
		cookie := &http.Cookie{
			Name:  cookies.CookieName,
			Value: c,
			Path:  "/",
		}
		request.AddCookie(cookie)
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	})
}
