package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set up Var that using
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method

		// Create log structure
		attributes := []slog.Attr{
			slog.String("method", method),
			slog.String("path", path),
			slog.String("query", query),
			slog.Int("status", status),
			slog.Duration("latency", latency),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				slog.Error("gin request error", slog.String("error", e))
			}
		}

		// Select level
		var level slog.Level
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		} else {
			level = slog.LevelInfo
		}

		slog.LogAttrs(c.Request.Context(), level, "gin request", attributes...)
	}
}
