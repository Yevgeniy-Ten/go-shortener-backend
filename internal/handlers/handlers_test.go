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

type TestHelper struct {
	Ctrl          *gomock.Controller
	Engine        *gin.Engine
	MockUserStore *handlers_mocks.MockuserStorage
	MockURLStore  *handlers_mocks.MockurlStorage
	Config        *handlers.Config
	Logger        logger.ZapLogger
	Storage       domain.Storage
}

func NewTestHelper(t *testing.T) *TestHelper {
	ctrl := gomock.NewController(t)
	mockUserStorage := handlers_mocks.NewMockuserStorage(ctrl)
	mockURLStorage := handlers_mocks.NewMockurlStorage(ctrl)

	cfg := &handlers.Config{
		ServerAddr:    "http://localhost:8080",
		DatabaseURL:   "soso",
		DatabaseError: false,
	}

	l := logger.ZapLogger{
		Log: zap.NewNop(),
	}

	storage := domain.Storage{
		User: mockUserStorage,
		URLS: mockURLStorage,
	}

	engine := handlers.InitHandlers(cfg, storage, &l)

	return &TestHelper{
		Ctrl:          ctrl,
		Engine:        engine,
		MockUserStore: mockUserStorage,
		MockURLStore:  mockURLStorage,
		Config:        cfg,
		Logger:        l,
		Storage:       storage,
	}
}

func TestInitHandlers(t *testing.T) {
	testHelper := NewTestHelper(t)

	defer testHelper.Ctrl.Finish()

	t.Run("Handlers inited", func(t *testing.T) {

		h := handlers.InitHandlers(testHelper.Config, testHelper.Storage, &testHelper.Logger)
		assert.NotNil(t, h)
	})
}
func TestCreateRouter(t *testing.T) {
	testHelper := NewTestHelper(t)
	defer testHelper.Ctrl.Finish()
	t.Run("Create routes created", func(t *testing.T) {
		h := handlers.NewHandler(testHelper.Config, testHelper.Storage, &testHelper.Logger)
		r := h.CreateRouter()
		assert.NotNil(t, r)
	})
}

func TestHandler_PostHandler(t *testing.T) {
	testHelper := NewTestHelper(t)
	defer testHelper.Ctrl.Finish()
	t.Run("PostHandler", func(t *testing.T) {
		testHelper.MockUserStore.EXPECT().Create().Return(1, nil)
		testHelper.MockURLStore.EXPECT().Save("http://practicum.yandex.ru", 1).Return("123", nil)
		request := httptest.NewRequest("POST", "/", strings.NewReader("http://practicum.yandex.ru"))
		recorder := httptest.NewRecorder()
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.Equal(t, "http://localhost:8080/123", recorder.Body.String())
	})
	t.Run("PostHandler with error", func(t *testing.T) {
		testHelper.MockUserStore.EXPECT().Create().Return(0, nil)
		request := httptest.NewRequest("POST", "/", strings.NewReader("practicum.yandex.ru"))
		recorder := httptest.NewRecorder()
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	})
	t.Run("PostHandler with error", func(t *testing.T) {
		testHelper.MockUserStore.EXPECT().Create().Return(1, nil)
		testHelper.MockURLStore.EXPECT().Save("http://practicum.yandex.ru", 1).Return("", errors.New("error"))
		request := httptest.NewRequest("POST", "/", strings.NewReader("http://practicum.yandex.ru"))
		recorder := httptest.NewRecorder()
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	})

}

func TestHandler_GetHandler(t *testing.T) {
	testHelper := NewTestHelper(t)
	defer testHelper.Ctrl.Finish()
	t.Run("GetHandler", func(t *testing.T) {
		testHelper.MockURLStore.EXPECT().GetURL("123").Return("http://practicum.yandex.ru", nil)
		request := httptest.NewRequest("GET", "/123", nil)
		recorder := httptest.NewRecorder()
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode)
		assert.Equal(t, "http://practicum.yandex.ru", result.Header.Get("Location"))
	})
	t.Run("GetHandler with error", func(t *testing.T) {
		testHelper.MockURLStore.EXPECT().GetURL("123").Return("", errors.New("error"))
		request := httptest.NewRequest("GET", "/123", nil)
		recorder := httptest.NewRecorder()
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusNotFound, result.StatusCode)
	})
}
func TestHandler_ShortenURLHandler(t *testing.T) {
	testHelper := NewTestHelper(t)
	defer testHelper.Ctrl.Finish()
	t.Run("ShortenURLHandler", func(t *testing.T) {

		testHelper.MockUserStore.EXPECT().Create().Return(1, nil)
		testHelper.MockURLStore.EXPECT().Save("http://practicum.yandex.ru", 1).Return("123", nil)
		request := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(`{"url":"http://practicum.yandex.ru"}`))
		recorder := httptest.NewRecorder()
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.Equal(t,
			"{\"result\":\"http://localhost:8080/123\"}", recorder.Body.String())
	})
	t.Run("ShortenURLHandler with error", func(t *testing.T) {
		testHelper.MockUserStore.EXPECT().Create().Return(0, nil)
		request := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(`{"url":"practicum.yandex.ru"}`))
		recorder := httptest.NewRecorder()
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	})
	t.Run("ShortenURLHandler with error", func(t *testing.T) {
		testHelper.MockUserStore.EXPECT().Create().Return(1, nil)
		testHelper.MockURLStore.EXPECT().Save("http://practicum.yandex.ru", 1).Return("", errors.New("error"))
		request := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(`{"url":"http://practicum.yandex.ru"}`))
		recorder := httptest.NewRecorder()
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	})
}

func TestHandler_GetUserUrls(t *testing.T) {
	testHelper := NewTestHelper(t)
	defer testHelper.Ctrl.Finish()
	t.Run("GetUserUrls", func(t *testing.T) {
		testHelper.MockUserStore.EXPECT().Create().Return(1, nil)
		testHelper.MockURLStore.EXPECT().GetUserURLs(
			1, "http://localhost:8080/").Return(
			[]domain.UserURLs{
				{OriginalURL: "http://localhost:8080/123", ShortURL: "123"},
			}, nil)
		request := httptest.NewRequest("GET", "/api/user/urls", nil)
		recorder := httptest.NewRecorder()
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.JSONEq(t, `[{"original_url":"http://localhost:8080/123","short_url":"123"}]`, recorder.Body.String())
	})
	t.Run("GetUserUrls with error", func(t *testing.T) {
		testHelper.MockUserStore.EXPECT().Create().Return(0, nil)
		testHelper.MockURLStore.EXPECT().GetUserURLs(
			0, "http://localhost:8080/").Return(
			[]domain.UserURLs{}, errors.New("error"))
		request := httptest.NewRequest("GET", "/api/user/urls", nil)
		recorder := httptest.NewRecorder()
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	})
	t.Run("GetUserUrls no urls", func(t *testing.T) {
		testHelper.MockUserStore.EXPECT().Create().Return(1, nil)
		testHelper.MockURLStore.EXPECT().GetUserURLs(
			1, "http://localhost:8080/").Return(
			[]domain.UserURLs{}, nil)
		request := httptest.NewRequest("GET", "/api/user/urls", nil)
		recorder := httptest.NewRecorder()
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusNoContent, result.StatusCode)
	})
}
func TestHandler_DeleteMyUrls(t *testing.T) {
	testHelper := NewTestHelper(t)
	defer testHelper.Ctrl.Finish()
	t.Run("DeleteMyUrls", func(t *testing.T) {
		ids := []string{"1", "2", "3"}
		testHelper.MockURLStore.EXPECT().DeleteURLs(ids, 1).Return(nil).AnyTimes()
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
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusAccepted, result.StatusCode)
	})
	t.Run("DeleteMyUrls with 401", func(t *testing.T) {
		request := httptest.NewRequest("DELETE", "/api/user/urls", nil)
		recorder := httptest.NewRecorder()
		testHelper.Engine.ServeHTTP(recorder, request)
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
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	})
}

func TestHandler_ShortenURLSHandler(t *testing.T) {
	testHelper := NewTestHelper(t)
	defer testHelper.Ctrl.Finish()
	t.Run("ShortenURLSHandler", func(t *testing.T) {
		testHelper.MockURLStore.EXPECT().SaveBatch(gomock.Any(), 1).Return(nil)
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
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusCreated, result.StatusCode)
		assert.JSONEq(t, `[{"correlation_id":"123","short_url":"http://localhost:8080/123"}]`, recorder.Body.String())
	})
	t.Run("ShortenURLSHandler with error", func(t *testing.T) {
		testHelper.MockURLStore.EXPECT().SaveBatch(gomock.Any(), 1).Return(errors.New("error"))
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
		testHelper.Engine.ServeHTTP(recorder, request)
		result := recorder.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	})
}
