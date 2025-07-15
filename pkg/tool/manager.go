package tool

import (
	"context"
	"fmt"
	"github.com/AdamShannag/api-mcp-server/pkg/request"
	"github.com/AdamShannag/api-mcp-server/pkg/resolver"
	"github.com/AdamShannag/api-mcp-server/pkg/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"log/slog"
)

type Option func(*Manager)

type Manager struct {
	executor    request.Executor
	argResolver resolver.ArgResolver
}

func NewManager(executor request.Executor, opts ...Option) *Manager {
	mgr := &Manager{
		executor:    executor,
		argResolver: resolver.NewDefaultTypeResolverRegistry(),
	}

	for _, opt := range opts {
		opt(mgr)
	}

	return mgr
}

func WithArgResolver(r resolver.ArgResolver) Option {
	return func(m *Manager) {
		m.argResolver = r
	}
}

func (tm *Manager) AddTool(mcpServer *server.MCPServer, tool types.Tool) {
	baseOptions := []mcp.ToolOption{
		mcp.WithDescription(tool.Description),
	}

	options := append(baseOptions, tm.toOptions(tool.Args)...)

	t := mcp.NewTool(tool.Name, options...)
	mcpServer.AddTool(t, tm.toolHandlerFactory(tool))

	slog.Debug("tool registered",
		slog.Group("tool",
			slog.String("name", tool.Name),
			slog.Int("args", len(tool.Args)),
		),
	)
}

func (tm *Manager) toOptions(args []types.Arg) []mcp.ToolOption {
	var options []mcp.ToolOption
	for _, arg := range args {
		options = append(options, tm.argResolver.ToToolOption(arg))
	}
	return options
}

func (tm *Manager) toolHandlerFactory(tool types.Tool) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := make(map[string]string)

		for _, arg := range tool.Args {
			val, err := tm.argResolver.Resolve(ctx, request, arg)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("invalid argument %q: %v", arg.Name, err)), nil
			}
			args[arg.Name] = val
		}

		resp, err := tm.executor.Execute(ctx, tool.Request, args)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("request failed: %v", err)), err
		}

		return mcp.NewToolResultText(resp), nil
	}
}
