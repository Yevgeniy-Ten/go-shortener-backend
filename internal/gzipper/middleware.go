package gzipper

import (
	"compress/gzip"
	"io"
	"net/http"
	logger "shorter/internal/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

type myWriter struct {
	writer *gzip.Writer
	gin.ResponseWriter
}

func (g *myWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

type myReader struct {
	reader *gzip.Reader
	io.ReadCloser
}

func (g *myReader) Read(p []byte) (n int, err error) {
	return g.reader.Read(p)
}

func RequestResponseGzipMiddleware(
	log logger.MyLogger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentEncodingHeader := c.GetHeader("Content-Encoding")
		if strings.Contains(contentEncodingHeader, "gzip") {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				log.Error("Ошибка чтения тела запроса")
				c.String(http.StatusBadRequest, "Ошибка чтения тела запроса")
				return
			}
			defer reader.Close()
			c.Request.Body = &myReader{
				reader:     reader,
				ReadCloser: c.Request.Body,
			}
		}
		acceptEncodingHeader := c.GetHeader("Accept-Encoding")
		if strings.Contains(acceptEncodingHeader, "gzip") {
			c.Writer.Header().Set("Content-Encoding", "gzip")
			writer := gzip.NewWriter(c.Writer)
			defer writer.Close()
			c.Writer = &myWriter{
				writer:         writer,
				ResponseWriter: c.Writer,
			}
		}
		c.Next()
	}
}
