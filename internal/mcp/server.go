package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AdamShannag/api-mcp-server/internal/auth"
	"github.com/AdamShannag/api-mcp-server/internal/middleware"
	"github.com/AdamShannag/api-mcp-server/internal/monitoring"
	"github.com/AdamShannag/api-mcp-server/pkg/tool"
	"github.com/AdamShannag/api-mcp-server/pkg/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	Version        = "0.2.0"
	serverName     = "API MCP Server"
	defaultSseHost = "127.0.0.1"
	defaultSsePort = 13080
)

type ServerOption func(*Server)

type Server struct {
	server        *server.MCPServer
	transport     string
	toolsFilePath string
	host          string
	port          string

	auth    *auth.Authenticator
	httpSrv *http.Server
}

func NewServer(transport string, opts ...ServerOption) *Server {
	s := &Server{
		transport: transport,
		host:      defaultSseHost,
		port:      strconv.Itoa(defaultSsePort),
	}

	for _, opt := range opts {
		opt(s)
	}

	options := []server.ServerOption{
		server.WithLogging(),
		server.WithRecovery(),
		server.WithToolCapabilities(true),
		server.WithHooks(s.getHooks()),
	}

	if s.auth != nil {
		options = append(options, server.WithToolHandlerMiddleware(s.auth.Middleware()))
		options = append(options, server.WithToolHandlerMiddleware(middleware.LoggingMiddleware))
	}

	s.server = server.NewMCPServer(
		serverName,
		Version,
		options...,
	)

	return s
}

func (s *Server) Run() error {
	if s.httpSrv != nil && s.transport != "stdio" {
		go func() {
			slog.Info("metrics server started", slog.String("addr", s.httpSrv.Addr))
			if err := s.httpSrv.ListenAndServe(); err != nil {
				slog.Error("metrics server error", slog.String("error", err.Error()))
			}
		}()
	}

	switch s.transport {
	case "sse":
		return s.runWithSSE()
	default:
		return server.ServeStdio(s.server)
	}
}

func (s *Server) LoadTools(manager *tool.Manager) error {
	data, err := os.ReadFile(s.toolsFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var tools []types.Tool
	decoder := json.NewDecoder(strings.NewReader(s.resolveEnvPlaceholders(string(data))))
	if err = decoder.Decode(&tools); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	for _, t := range tools {
		manager.AddTool(s.server, t)
	}

	slog.Info("tools loaded", slog.Int("count", len(tools)))
	return nil
}

func WithAuth(a *auth.Authenticator) ServerOption {
	return func(s *Server) {
		s.auth = a
	}
}

func WithToolsFile(path string) ServerOption {
	return func(s *Server) {
		s.toolsFilePath = path
	}
}

func WithHost(host string) ServerOption {
	return func(s *Server) {
		if host == "" {
			return
		}
		s.host = host
	}
}

func WithPort(port string) ServerOption {
	return func(s *Server) {
		if port == "" {
			return
		}
		s.port = port
	}
}

func WithHttpServer(server *http.Server) ServerOption {
	return func(s *Server) {
		s.httpSrv = server
	}
}

func (s *Server) getHooks() *server.Hooks {
	hooks := &server.Hooks{}

	hooks.AddOnRegisterSession(func(ctx context.Context, session server.ClientSession) {
		slog.Info("client connected", slog.String("sessionId", session.SessionID()))
		monitoring.SessionStarts.Inc()
		monitoring.ActiveSessions.Inc()
	})

	hooks.AddOnUnregisterSession(func(ctx context.Context, session server.ClientSession) {
		slog.Warn("client disconnected", slog.String("sessionId", session.SessionID()))
		monitoring.SessionCloses.Inc()
		monitoring.ActiveSessions.Dec()
	})

	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		slog.Info("processing request", slog.String("method", string(method)))
	})

	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		slog.Error("error occurred", slog.String("method", string(method)), slog.String("error", err.Error()))
		monitoring.ErrorsTotal.WithLabelValues(string(method)).Inc()
	})

	return hooks
}

func (s *Server) runWithSSE() error {
	sseServer := server.NewSSEServer(s.server,
		server.WithBaseURL(fmt.Sprintf("http://:%s", s.port)),
		server.WithSSEContextFunc(s.auth.FromRequest),
	)

	var runErr error

	s.startWithGracefulShutdown(
		func() {
			slog.Info("sse server started", slog.String("host", s.host), slog.String("port", s.port))
			if err := sseServer.Start(s.host + ":" + s.port); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error("server error", slog.String("error", err.Error()))
				runErr = err
			}
		},
		func(ctx context.Context) error {
			var g errgroup.Group
			g.Go(func() error { return sseServer.Shutdown(ctx) })
			if s.httpSrv != nil {
				g.Go(func() error { return s.httpSrv.Shutdown(ctx) })
			}
			return g.Wait()
		},
	)

	return runErr
}

func (s *Server) startWithGracefulShutdown(initFunc func(), shutdownFunc func(context.Context) error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go initFunc()

	<-sigChan
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := shutdownFunc(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server stopped")
}

// resolveEnvPlaceholders resolves placeholders like {{env VAR_NAME:default}} with environment values.
func (s *Server) resolveEnvPlaceholders(in string) string {
	var b strings.Builder
	b.Grow(len(in))

	i := 0
	for i < len(in) {
		start := strings.Index(in[i:], "{{env ")
		if start == -1 {
			b.WriteString(in[i:])
			break
		}
		start += i
		b.WriteString(in[i:start])
		end := strings.Index(in[start:], "}}")
		if end == -1 {
			b.WriteString(in[start:])
			break
		}
		end += start

		content := strings.TrimSpace(in[start+6 : end])
		varName := strings.TrimSpace(content)
		defaultValue := ""
		if idx := strings.Index(content, ":"); idx != -1 {
			varName = strings.TrimSpace(content[:idx])
			defaultValue = strings.TrimSpace(content[idx+1:])
		}
		if val := os.Getenv(varName); val != "" {
			b.WriteString(val)
		} else {
			b.WriteString(defaultValue)
		}
		i = end + 2
	}
	return b.String()
}
