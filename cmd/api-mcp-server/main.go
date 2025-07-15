package main

import (
	"flag"
	"github.com/AdamShannag/api-mcp-server/internal/auth"
	"github.com/AdamShannag/api-mcp-server/internal/mcp"
	"github.com/AdamShannag/api-mcp-server/pkg/request"
	"github.com/AdamShannag/api-mcp-server/pkg/tool"
	"log"
	"net/http"
	"os"
	"time"
)

const httpClientTimeout = 30 * time.Second

func main() {
	var transport string
	var toolsFilePath string
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or sse)")

	flag.StringVar(&toolsFilePath, "c", "./config.json", "Tools config file path")
	flag.StringVar(&toolsFilePath, "config", "./config.json", "Tools config file path")
	flag.Parse()

	manager := tool.NewManager(request.NewExecutor(
		request.WithHttpClient(&http.Client{Timeout: httpClientTimeout}),
	))

	s := mcp.NewServer(transport,
		mcp.WithHost(os.Getenv("API_MCP_HOST")),
		mcp.WithPort(os.Getenv("API_MCP_PORT")),
		mcp.WithToolsFile(toolsFilePath),
		mcp.WithAuth(auth.NewAuthenticator("sse", os.Getenv("API_MCP_SSE_API_KEY"))),
	)

	err := s.LoadTools(manager)

	if err != nil {
		log.Fatal(err)
	}

	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
