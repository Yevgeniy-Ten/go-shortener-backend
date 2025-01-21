package cookies

import (
	"fmt"
	"net/http"
	"shorter/internal/domain"
	"shorter/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

const (
	CookieName = "token"
	MaxAge     = 3600
)

type SomeRepo interface {
	Create() (int, error)
}

const hashKey = []byte("your-secret-hash-key") // 16 bytes or more
const blockKey = []byte("your-block-key-16byt")

var s = securecookie.New(hashKey, blockKey)

func CreateUserMiddleware(l *logger.ZapLogger, repo SomeRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
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
			encoded, err := s.Encode(CookieName, userID)
			if err != nil {
				fmt.Println(err.Error())
				c.JSON(http.StatusInternalServerError, domain.ResponseError{
					Description: "Error encoding cookie",
				})
				c.Abort()
				return
			}
			c.SetCookie(CookieName, encoded, MaxAge, "/", "", false, true)
			c.Status(http.StatusUnauthorized)
			c.Abort()

		}
		c.Next()
	}
}
