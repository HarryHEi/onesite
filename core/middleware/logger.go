package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"onesite/common/log"
)

func Logger() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}
		statusCode := c.Writer.Status()

		logFields := []zap.Field{
			zap.Time("TimeStamp", time.Now()),
			zap.Duration("Latency", time.Now().Sub(start)),
			zap.String("ClientIP", c.ClientIP()),
			zap.String("Method", c.Request.Method),
			zap.Int("StatusCode", statusCode),
			zap.String("ErrorMessage", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Int("BodySize", c.Writer.Size()),
			zap.String("Path", path),
		}

		switch {
		case statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices:
			log.Info("Gin", logFields...)
		case statusCode >= http.StatusMultipleChoices && statusCode < http.StatusBadRequest:
			log.Warn("Gin", logFields...)
		case statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError:
			log.Error("Gin", logFields...)
		default:
			log.Warn("Gin", logFields...)
		}
	}
}
