package config

import (
	"log/slog"
)

var logger = slog.Default()

// Logger returns the logger instance.
func Logger() *slog.Logger {
	return logger
}

// SetLogger sets the default logger's level.
func SetLogger(l *slog.Logger) {
	logger = l
}
