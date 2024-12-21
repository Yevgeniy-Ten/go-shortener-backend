package gzipper

import (
	"compress/gzip"
	"io"
	"net/http"
	"shorter/internal/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) Close() error {
	return g.writer.Close()
}
func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

type gzipReader struct {
	reader *gzip.Reader
	io.ReadCloser
}

func (g *gzipReader) Read(p []byte) (n int, err error) {
	return g.reader.Read(p)
}

func (g *gzipReader) Close() error {
	return g.reader.Close()
}
func RequestResponseGzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentEncodingHeader := c.GetHeader("Content-Encoding")
		if strings.Contains(contentEncodingHeader, "gzip") {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				logger.Log.Error("Ошибка чтения тела запроса")
				c.String(http.StatusBadRequest, "Ошибка чтения тела запроса")
				return
			}
			defer reader.Close()
			c.Request.Body = &gzipReader{
				reader:     reader,
				ReadCloser: c.Request.Body,
			}
		}
		//acceptEncodingHeader := c.GetHeader("Accept-Encoding")
		//if strings.Contains(acceptEncodingHeader, "gzip") {
		//	c.Writer.Header().Set("Content-Encoding", "gzip")
		//	writer := gzip.NewWriter(c.Writer)
		//	c.Writer = &gzipWriter{
		//		ResponseWriter: c.Writer,
		//		writer:         writer,
		//	}
		//	defer writer.Close()
		//	c.Next()
		//	return
		//}
		c.Next()
	}
}
