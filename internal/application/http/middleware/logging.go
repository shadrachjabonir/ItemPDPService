package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LoggingMiddleware creates a logging middleware
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.Info().
			Str("method", param.Method).
			Str("path", param.Path).
			Int("status", param.StatusCode).
			Dur("latency", param.Latency).
			Str("client_ip", param.ClientIP).
			Str("user_agent", param.Request.UserAgent()).
			Int("body_size", param.BodySize).
			Msg("HTTP Request")
		return ""
	})
}

// StructuredLoggingMiddleware provides structured logging using zerolog
func StructuredLoggingMiddleware(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get client IP
		clientIP := c.ClientIP()

		// Get status code
		statusCode := c.Writer.Status()

		// Get body size
		bodySize := c.Writer.Size()

		// Get user agent
		userAgent := c.Request.UserAgent()

		// Build log entry
		logEvent := logger.Info()

		if len(c.Errors) > 0 {
			logEvent = logger.Error().Strs("errors", c.Errors.Errors())
		}

		if raw != "" {
			path = path + "?" + raw
		}

		logEvent.
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("client_ip", clientIP).
			Str("user_agent", userAgent).
			Int("body_size", bodySize).
			Msg("HTTP Request")
	}
} 