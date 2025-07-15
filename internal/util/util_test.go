package util_test

import (
	"github.com/AdamShannag/api-mcp-server/internal/util"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{"debug lowercase", "debug", slog.LevelDebug},
		{"DEBUG uppercase", "DEBUG", slog.LevelDebug},
		{"info lowercase", "info", slog.LevelInfo},
		{"INFO uppercase", "INFO", slog.LevelInfo},
		{"warn lowercase", "warn", slog.LevelWarn},
		{"WARNING uppercase", "WARNING", slog.LevelWarn},
		{"error lowercase", "error", slog.LevelError},
		{"ERROR uppercase", "ERROR", slog.LevelError},
		{"unknown value", "trace", slog.LevelInfo},
		{"empty string", "", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.GetLogLevel(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
