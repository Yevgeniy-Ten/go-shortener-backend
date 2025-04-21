package cookies

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserFromCookie(t *testing.T) {
	t.Run("validCookie", func(t *testing.T) {
		userID := 1
		encoded, err := s.Encode(CookieName, userID)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Request.Header.Set("Cookie", fmt.Sprintf("%s=%s", CookieName, encoded))
		userID, err = GetUserFromCookie(c)
		assert.NoError(t, err)
		assert.Equal(t, 1, userID)
	})
	t.Run("invalidCookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		userID, err := GetUserFromCookie(c)
		assert.Error(t, err)
		assert.Equal(t, 0, userID)
	})
}

func TestCreateCookie(t *testing.T) {
	t.Run("validCookie", func(t *testing.T) {
		userID := 1
		encoded, err := CreateCookie(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, encoded)
	})
}
