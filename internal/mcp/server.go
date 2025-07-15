package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AdamShannag/api-mcp-server/internal/auth"
	"github.com/AdamShannag/api-mcp-server/pkg/tool"
	"github.com/AdamShannag/api-mcp-server/pkg/types"
	"github.com/mark3labs/mcp-go/server"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	serverName     = "API MCP Server"
	version        = "0.0.1"
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

	auth *auth.Authenticator
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
	}

	if s.auth != nil {
		options = append(options, server.WithToolHandlerMiddleware(s.auth.Middleware()))
	}

	s.server = server.NewMCPServer(
		serverName,
		version,
		options...,
	)

	return s
}

func (s *Server) Run() error {
	switch s.transport {
	case "sse":
		sseServer := server.NewSSEServer(s.server,
			server.WithBaseURL(fmt.Sprintf("http://:%s", s.port)),
			server.WithSSEContextFunc(s.auth.FromRequest),
		)

		s.startWithGracefulShutdown(func() {
			log.Println("SSE server listening on " + s.host + ":" + s.port)
			if err := sseServer.Start(s.host + ":" + s.port); err != nil {
				log.Fatalf("Server error: %v", err)
			}
		},
			func(ctx context.Context) error {
				return sseServer.Shutdown(ctx)
			})

	default:
		if err := server.ServeStdio(s.server); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}

	return nil
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
		s.host = host
	}
}

func WithPort(port string) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) startWithGracefulShutdown(initFunc func(), shutdownFunc func(context.Context) error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go initFunc()

	<-sigChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := shutdownFunc(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	log.Println("Server stopped")
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
