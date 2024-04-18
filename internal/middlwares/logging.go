package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ResponseWriterInterceptor struct {
	gin.ResponseWriter

	code int
}

func WrapperResponse(w gin.ResponseWriter) *ResponseWriterInterceptor {
	return &ResponseWriterInterceptor{ResponseWriter: w}
}

func (rw *ResponseWriterInterceptor) StatusCode() int {
	if rw.code == 0 {
		return http.StatusOK
	}

	return rw.code
}

func (rw *ResponseWriterInterceptor) WriteHeader(statusCode int) {
	rw.code = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		wrapper := WrapperResponse(c.Writer)
		c.Writer = wrapper
		start := time.Now()
		c.Next()

		logger.WithContext(c).WithFields(
			logrus.Fields{
				"method":   c.Request.Method,
				"host":     c.Request.Host,
				"url":      c.Request.URL.EscapedPath(),
				"status":   wrapper.StatusCode(),
				"duration": time.Since(start),
			},
		).Info("HTTP Request")
	}
}
