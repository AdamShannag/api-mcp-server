package util

import (
	"log/slog"
	"strings"
)

func GetLogLevel(lv string) slog.Level {
	switch strings.ToUpper(lv) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
