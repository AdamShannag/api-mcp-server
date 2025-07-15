package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	assert.Contains(t, err.Error(), "failed to decode JSON")
}

func TestServer_LoadTools_FileMissing(t *testing.T) {
	s := NewServer("stdio", WithToolsFile("/not/found/tools.json"))
	manager := &tool.Manager{}

	err := s.LoadTools(manager)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")
}

func TestResolveEnvPlaceholders(t *testing.T) {
	s := &Server{}

	t.Setenv("HOST_URL", "https://example.com")
	t.Setenv("API_KEY", "secret_key")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "env var present, no default",
			input:    "Connect to {{env HOST_URL}} now",
			expected: "Connect to https://example.com now",
		},
		{
			name:     "env var missing, default used",
			input:    "API key is {{env MISSING_KEY:default_key}}",
			expected: "API key is default_key",
		},
		{
			name:     "env var missing, no default",
			input:    "Value: {{env MISSING}}",
			expected: "Value: ",
		},
		{
			name:     "multiple placeholders",
			input:    "Host: {{env HOST_URL}}, Key: {{env API_KEY:default_key}}",
			expected: "Host: https://example.com, Key: secret_key",
		},
		{
			name:     "no placeholders",
			input:    "Just a regular string",
			expected: "Just a regular string",
		},
		{
			name:     "malformed placeholder no closing braces",
			input:    "Value {{env HOST_URL",
			expected: "Value {{env HOST_URL",
		},
		{
			name:     "placeholder with spaces and default",
			input:    "URL: {{env   HOST_URL  :   https://default.com  }}",
			expected: "URL: https://example.com",
		},
		{
			name:     "placeholder with empty default",
			input:    "Empty default {{env MISSING:}} end",
			expected: "Empty default  end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := s.resolveEnvPlaceholders(tt.input)
			assert.Equal(t, tt.expected, out)
		})
	}
}

func BenchmarkResolveEnvPlaceholders_Manual(b *testing.B) {
	s := &Server{}
	b.Setenv("API_KEY", "live_key_123")
	tmpl := largeTemplate()

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = s.resolveEnvPlaceholders(tmpl)
	}
}

func largeTemplate() string {
	var sb strings.Builder
	sb.Grow(1_000_000)

	for i := 0; i < 1000; i++ {
		sb.WriteString(`{
			"name": "Tool ` + string(rune('A'+(i%26))) + `",
			"request": {
				"host": "{{env API_KEY:default_key}}",
				"headers": {
					"Authorization": "Bearer {{env API_KEY}}"
				}
			}
		},`)
	}
	return sb.String()
}
