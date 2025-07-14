package auth

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"net/http"
	"strings"
)

type contextKey string

const authContextKey = contextKey("auth-key")

type Authenticator struct {
	ApiKey    string
	Transport string
}

func NewAuthenticator(transport string, token string) *Authenticator {
	return &Authenticator{
		ApiKey:    token,
		Transport: transport,
	}
}

func (a *Authenticator) FromRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, authContextKey, r.Header.Get("Authorization"))
}

func (a *Authenticator) Middleware() server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			switch a.Transport {
			case "stdio":
				return next(ctx, req)

			case "sse":
				ok, err := a.authenticate(ctx)
				if err != nil {
					return nil, fmt.Errorf("authentication error: %w", err)
				}
				if !ok {
					return nil, errors.New("unauthorized request")
				}
				return next(ctx, req)

			default:
				return nil, fmt.Errorf("unknown transport type: %s", a.Transport)
			}
		}
	}
}

func (a *Authenticator) authenticate(ctx context.Context) (bool, error) {
	if a.ApiKey == "" {
		return true, nil
	}

	rawHeader, ok := ctx.Value(authContextKey).(string)
	if !ok || rawHeader == "" {
		return false, errors.New("missing Authorization header")
	}

	token := strings.TrimPrefix(rawHeader, "Bearer ")
	if token == "" {
		return false, errors.New("empty token in Authorization header")
	}

	if subtle.ConstantTimeCompare([]byte(a.ApiKey), []byte(token)) != 1 {
		return false, errors.New("invalid auth token")
	}

	return true, nil
}
