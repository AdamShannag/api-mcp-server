package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/AdamShannag/api-mcp-server/internal/auth"
	"github.com/AdamShannag/api-mcp-server/pkg/tool"
	"github.com/AdamShannag/api-mcp-server/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewServer_WithDefaults(t *testing.T) {
	s := NewServer("sse")

	assert.Equal(t, "sse", s.transport)
	assert.Equal(t, defaultSseHost, s.host)
	assert.Equal(t, strconv.Itoa(defaultSsePort), s.port)
	assert.Nil(t, s.auth)
	assert.Empty(t, s.toolsFilePath)
}

func TestNewServer_WithOptions(t *testing.T) {
	authenticator := auth.NewAuthenticator("sse", "super-secret")

	s := NewServer("sse",
		WithAuth(authenticator),
		WithToolsFile("/tmp/tools.json"),
		WithHost("0.0.0.0"),
		WithPort("9999"),
	)

	assert.Equal(t, "sse", s.transport)
	assert.Equal(t, "0.0.0.0", s.host)
	assert.Equal(t, "9999", s.port)
	assert.Equal(t, "/tmp/tools.json", s.toolsFilePath)
	assert.Equal(t, authenticator, s.auth)
}

func TestServer_LoadTools_Success(t *testing.T) {
	tmpDir := t.TempDir()
	toolsFile := filepath.Join(tmpDir, "tools.json")

	mockTools := []types.Tool{
		{
			Name:        "Ping",
			Description: "Ping a service",
			Request: types.Request{
				Host:     "example.com",
				Endpoint: "/ping",
				Method:   "GET",
				Secure:   true,
			},
		},
	}
	data, _ := json.Marshal(mockTools)
	_ = os.WriteFile(toolsFile, data, 0644)

	s := NewServer("stdio", WithToolsFile(toolsFile))
	manager := &tool.Manager{}

	err := s.LoadTools(manager)

	assert.NoError(t, err)
}

func TestServer_LoadTools_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	badJSON := filepath.Join(tmpDir, "bad.json")

	_ = os.WriteFile(badJSON, []byte("{invalid-json}"), 0644)

	s := NewServer("stdio", WithToolsFile(badJSON))
	manager := &tool.Manager{}

	err := s.LoadTools(manager)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal")
}

func TestServer_LoadTools_FileMissing(t *testing.T) {
	s := NewServer("stdio", WithToolsFile("/not/found/tools.json"))
	manager := &tool.Manager{}

	err := s.LoadTools(manager)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")
}
