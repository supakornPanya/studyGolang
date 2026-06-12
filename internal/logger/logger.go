package logger

import (
	"log/slog"
	"os"
)

func InitLogger() {
	var handler slog.Handler
	mode := os.Getenv("APP_ENV")
	if mode == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}