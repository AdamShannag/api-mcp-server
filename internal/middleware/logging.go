package middleware

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"log/slog"
	"strings"
	"time"
)

func LoggingMiddleware(next server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		start := time.Now()
		sessionID := server.ClientSessionFromContext(ctx).SessionID()
		toolName := req.Params.Name

		argPairs := make([]string, 0, len(req.GetArguments()))
		for k, v := range req.GetArguments() {
			argPairs = append(argPairs, fmt.Sprintf("%s=%s", k, v))
		}

		slog.Info("tool call started",
			slog.String("tool", toolName),
			slog.String("sessionId", sessionID),
			slog.String("args", strings.Join(argPairs, ", ")),
		)

		result, err := next(ctx, req)
		duration := time.Since(start)

		if err != nil {
			slog.Error("tool call failed",
				slog.String("tool", toolName),
				slog.String("sessionId", sessionID),
				slog.Duration("duration", duration),
				slog.String("error", err.Error()),
			)
			return result, err
		}

		slog.Info("tool call completed",
			slog.String("tool", toolName),
			slog.String("sessionId", sessionID),
			slog.Duration("duration", duration),
		)

		return result, nil
	}
}
