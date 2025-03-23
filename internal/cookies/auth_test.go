package cookies_test

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"shorter/internal/cookies"
	"shorter/internal/cookies/mocks"
	"shorter/internal/logger"
	"testing"
)

func handlerStub(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
func TestCreateUserMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	userID := 1
	mockLogger := &logger.ZapLogger{}
	t.Run("should create user when withDatabase is true", func(t *testing.T) {
		mockUserRepo.EXPECT().Create().Return(userID, nil)
		router := gin.New()
		router.Use(cookies.CreateUserMiddleware(true, mockLogger, mockUserRepo))
		router.GET("/", handlerStub)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("should skip middleware when withDatabase is false", func(t *testing.T) {
		router := gin.New()
		router.Use(cookies.CreateUserMiddleware(false, mockLogger, mockUserRepo))
		router.GET("/", handlerStub)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("get cookie from response", func(t *testing.T) {
		mockUserRepo.EXPECT().Create().Return(userID, nil)
		router := gin.New()
		router.Use(cookies.CreateUserMiddleware(true, mockLogger, mockUserRepo))
		router.GET("/", handlerStub)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		cookie := w.Header().Get("Set-Cookie")
		assert.Contains(t, cookie, "token")
	})
}
