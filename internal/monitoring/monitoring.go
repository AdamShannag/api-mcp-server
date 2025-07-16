package monitoring

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const namespace = "api_mcp_server"

var (
	ToolInvocations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "tool_invocations_total",
			Help:      "Total number of tool invocations",
		},
		[]string{"tool"},
	)

	ToolLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "tool_duration_seconds",
			Help:      "Tool execution duration in seconds",
			Buckets:   prometheus.ExponentialBuckets(0.01, 2, 12),
		},
		[]string{"tool"},
	)

	SessionStarts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "session_starts_total",
			Help:      "Total number of sessions started",
		},
	)

	SessionCloses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "session_closes_total",
			Help:      "Total number of sessions closed",
		},
	)

	ActiveSessions = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "active_sessions",
			Help:      "Current number of active sessions",
		},
	)

	ErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "errors_total",
			Help:      "Total number of errors encountered per method",
		},
		[]string{"method"},
	)
)

func NewHttpServer(enabled bool, port string) *http.Server {
	if !enabled {
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(newRegistry(), promhttp.HandlerOpts{}))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
}

func newRegistry() *prometheus.Registry {
	reg := prometheus.NewRegistry()

	reg.MustRegister(
		ToolInvocations,
		ToolLatency,
		SessionStarts,
		SessionCloses,
		ActiveSessions,
		ErrorsTotal,
	)

	return reg
}
