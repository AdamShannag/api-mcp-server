package tool

import (
	"context"
	"errors"
	"github.com/AdamShannag/api-mcp-server/pkg/resolver"
	"github.com/AdamShannag/api-mcp-server/pkg/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManager_ToolHandlerFactory_Success(t *testing.T) {
	mockExec := &mockExecutor{
		output: `{"result":"ok"}`,
		err:    nil,
	}

	mockRes := &mockResolver{
		resolveFunc: func(ctx context.Context, req resolver.CallToolRequest, arg types.Arg) (string, error) {
			if arg.Name == "id" {
				return "42", nil
			}
			return "", errors.New("unexpected arg")
		},
		toToolOptionFn: func(arg types.Arg) mcp.ToolOption {
			return mcp.WithString(arg.Name)
		},
	}

	mgr := NewManager(mockExec, WithArgResolver(mockRes))

	toolDef := types.Tool{
		Name:        "GetTodo",
		Description: "Get todo by ID",
		Request: types.Request{
			Host:       "example.com",
			Endpoint:   "/todos/:id",
			Method:     "GET",
			Secure:     false,
			PathParams: []string{"id"},
		},
		Args: []types.Arg{
			{Name: "id", Type: "string", Required: true},
		},
	}

	handler := mgr.toolHandlerFactory(toolDef)

	resp, err := handler(context.Background(), mcp.CallToolRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.IsError)
	assert.Contains(t, resp.Content[0].(mcp.TextContent).Text, `"result":"ok"`)
}

func TestManager_ToolHandlerFactory_ArgResolveError(t *testing.T) {
	mockExec := &mockExecutor{}

	mockRes := &mockResolver{
		resolveFunc: func(ctx context.Context, req resolver.CallToolRequest, arg types.Arg) (string, error) {
			return "", errors.New("bad argument")
		},
		toToolOptionFn: func(arg types.Arg) mcp.ToolOption {
			return mcp.WithString(arg.Name)
		},
	}

	mgr := NewManager(mockExec, WithArgResolver(mockRes))

	toolDef := types.Tool{
		Name: "BrokenTool",
		Args: []types.Arg{
			{Name: "fail", Type: "string", Required: true},
		},
	}

	handler := mgr.toolHandlerFactory(toolDef)

	resp, err := handler(context.Background(), mcp.CallToolRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.IsError)
	assert.Contains(t, resp.Content[0].(mcp.TextContent).Text, "invalid argument")
}

func TestManager_ToolHandlerFactory_ExecuteError(t *testing.T) {
	mockExec := &mockExecutor{
		err: errors.New("execution failed"),
	}

	mockRes := &mockResolver{
		resolveFunc: func(ctx context.Context, req resolver.CallToolRequest, arg types.Arg) (string, error) {
			return "42", nil
		},
		toToolOptionFn: func(arg types.Arg) mcp.ToolOption {
			return mcp.WithString(arg.Name)
		},
	}

	mgr := NewManager(mockExec, WithArgResolver(mockRes))

	toolDef := types.Tool{
		Name: "ExecFailTool",
		Args: []types.Arg{
			{Name: "id", Type: "string", Required: true},
		},
		Request: types.Request{
			Host:       "example.com",
			Endpoint:   "/execfail/:id",
			Method:     "GET",
			PathParams: []string{"id"},
		},
	}

	handler := mgr.toolHandlerFactory(toolDef)

	resp, err := handler(context.Background(), mcp.CallToolRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.IsError)
	assert.Contains(t, resp.Content[0].(mcp.TextContent).Text, "request failed")
}

type mockResolver struct {
	resolveFunc    func(ctx context.Context, req resolver.CallToolRequest, arg types.Arg) (string, error)
	toToolOptionFn func(arg types.Arg) mcp.ToolOption
}

func (m *mockResolver) Resolve(ctx context.Context, req resolver.CallToolRequest, arg types.Arg) (string, error) {
	if m.resolveFunc != nil {
		return m.resolveFunc(ctx, req, arg)
	}
	return "", nil
}

func (m *mockResolver) ToToolOption(arg types.Arg) mcp.ToolOption {
	if m.toToolOptionFn != nil {
		return m.toToolOptionFn(arg)
	}
	return nil
}

type mockExecutor struct {
	output string
	err    error
}

func (m *mockExecutor) Execute(_ context.Context, _ types.Request, _ map[string]string) (string, error) {
	return m.output, m.err
}
