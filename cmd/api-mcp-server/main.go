package main

import (
	"flag"
	"fmt"
	"github.com/AdamShannag/api-mcp-server/internal/auth"
	"github.com/AdamShannag/api-mcp-server/internal/mcp"
	"github.com/AdamShannag/api-mcp-server/internal/monitoring"
	"github.com/AdamShannag/api-mcp-server/internal/util"
	"github.com/AdamShannag/api-mcp-server/pkg/request"
	"github.com/AdamShannag/api-mcp-server/pkg/tool"
	"github.com/lmittmann/tint"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const httpClientTimeout = 30 * time.Second

func main() {
	var (
		transport     string
		toolsFilePath string
		showVersion   bool
		enableMetrics bool
		metricsPort   string
	)
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or sse)")

	flag.StringVar(&toolsFilePath, "c", "./config.json", "Tools config file path")
	flag.StringVar(&toolsFilePath, "config", "./config.json", "Tools config file path")

	flag.BoolVar(&showVersion, "v", false, "Show version and exit")
	flag.BoolVar(&showVersion, "version", false, "Show version and exit")

	flag.BoolVar(&enableMetrics, "m", false, "Start metrics server")
	flag.BoolVar(&enableMetrics, "metrics", false, "Start metrics server")
	flag.StringVar(&metricsPort, "metrics-port", "8080", "Port for metrics endpoint (default: 8080)")
	flag.Parse()

	if showVersion {
		fmt.Println("api-mcp-server", mcp.Version)
		return
	}

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      util.GetLogLevel(os.Getenv("LOG_LEVEL")),
			TimeFormat: time.DateTime,
		}),
	))

	manager := tool.NewManager(request.NewExecutor(
		request.WithHttpClient(&http.Client{Timeout: httpClientTimeout}),
	))

	s := mcp.NewServer(transport,
		mcp.WithHost(os.Getenv("API_MCP_HOST")),
		mcp.WithPort(os.Getenv("API_MCP_PORT")),
		mcp.WithToolsFile(toolsFilePath),
		mcp.WithAuth(auth.NewAuthenticator("sse", os.Getenv("API_MCP_SSE_API_KEY"))),
		mcp.WithHttpServer(monitoring.NewHttpServer(enableMetrics, metricsPort)),
	)

	err := s.LoadTools(manager)

	if err != nil {
		log.Fatal(err)
	}

	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
