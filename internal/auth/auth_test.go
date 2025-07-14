package auth

import (
	"context"
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

func dummyHandler(called *bool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		*called = true
		return &mcp.CallToolResult{
			Result:  mcp.Result{},
			Content: nil,
			IsError: false,
		}, nil
	}
}

func TestAuthenticator_StdioSkipsAuth(t *testing.T) {
	auth := NewAuthenticator("stdio", "any-token")

	called := false
	middleware := auth.Middleware()
	handler := middleware(dummyHandler(&called))

	_, err := handler(context.Background(), mcp.CallToolRequest{})
	assert.NoError(t, err)
	assert.True(t, called, "stdio should skip auth and call next")
}

func TestAuthenticator_SuccessfulAuth(t *testing.T) {
	auth := NewAuthenticator("sse", "my-token")

	ctx := context.WithValue(context.Background(), authContextKey, "Bearer my-token")

	called := false
	handler := auth.Middleware()(dummyHandler(&called))

	_, err := handler(ctx, mcp.CallToolRequest{})
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestAuthenticator_InvalidToken(t *testing.T) {
	auth := NewAuthenticator("sse", "expected-token")

	ctx := context.WithValue(context.Background(), authContextKey, "Bearer wrong-token")

	called := false
	handler := auth.Middleware()(dummyHandler(&called))

	_, err := handler(ctx, mcp.CallToolRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid auth token")
	assert.False(t, called)
}

func TestAuthenticator_MissingAuthHeader(t *testing.T) {
	auth := NewAuthenticator("sse", "secret")

	ctx := context.WithValue(context.Background(), authContextKey, "")

	handler := auth.Middleware()(dummyHandler(nil))
	_, err := handler(ctx, mcp.CallToolRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing Authorization header")
}

func TestAuthenticator_EmptyBearer(t *testing.T) {
	auth := NewAuthenticator("sse", "secret")

	ctx := context.WithValue(context.Background(), authContextKey, "Bearer ")

	handler := auth.Middleware()(dummyHandler(nil))
	_, err := handler(ctx, mcp.CallToolRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty token")
}

func TestAuthenticator_UnknownTransport(t *testing.T) {
	auth := NewAuthenticator("weird", "key")

	handler := auth.Middleware()(dummyHandler(nil))
	_, err := handler(context.Background(), mcp.CallToolRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown transport type")
}

func TestFromRequest_ExtractsHeader(t *testing.T) {
	auth := NewAuthenticator("sse", "any")

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer test123")

	ctx := auth.FromRequest(context.Background(), req)

	val := ctx.Value(authContextKey)
	assert.Equal(t, "Bearer test123", val)
}
