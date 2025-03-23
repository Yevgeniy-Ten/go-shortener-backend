package cookies

import (
	"errors"
	"fmt"
	"net/http"
	"shorter/internal/domain"
	"shorter/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

// mockgen -source=internal/cookies/auth.go -destination=internal/cookies/mocks/mock_user_service.go -package=mocks
// Variables for cookies
const (
	CookieName = "token"
	MaxAge     = 3600
)

// UserRepo interface for creating user
//
//go:generate mockgen -source=internal/cookies/auth.go -destination=internal/cookies/mocks/mock_user_service.go -package=mocks
type UserRepo interface {
	Create() (int, error)
}

var hashKey = []byte("my-secret-hash-key") // 16 bytes or more

var s = securecookie.New(hashKey, nil)

// CreateUserMiddleware creates a user middleware
func CreateUserMiddleware(withDatabase bool, l *logger.ZapLogger, repo UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !withDatabase {
			c.Next()
			return
		}
		userCookie, err := c.Cookie(CookieName)
		if err != nil || userCookie == "" {
			userID, err := repo.Create()
			if err != nil {
				l.Log.Error("middleware: Error creating user")
				c.JSON(http.StatusInternalServerError, domain.ResponseError{
					Description: "Error creating user",
				})
				c.Abort()
				return
			}
			encoded, err := CreateCookie(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, domain.ResponseError{
					Description: "Error encoding cookie",
				})
				c.Abort()
				return
			}
			c.SetCookie(CookieName, encoded, MaxAge, "/", "", false, false)
			c.Request.Header.Set("Cookie", fmt.Sprintf("%s=%s", CookieName, encoded))
		}
		c.Next()
	}
}

// CreateCookie creates a cookie
func CreateCookie(userID int) (string, error) {
	encoded, err := s.Encode(CookieName, userID)
	if err != nil {
		return "", errors.New("error encoding cookie")
	}
	return encoded, nil
}

// GetUserFromCookie gets user from cookie
func GetUserFromCookie(c *gin.Context) (int, error) {
	userCookie, err := c.Cookie(CookieName)

	if err != nil {
		return 0, errors.New("no cookie")
	}

	var userID int
	if err := s.Decode(CookieName, userCookie, &userID); err != nil {
		return 0, errors.New("error decoding cookie")
	}
	return userID, nil
}
