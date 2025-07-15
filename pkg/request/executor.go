package request

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AdamShannag/api-mcp-server/pkg/types"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

type Option func(*executor)

type Executor interface {
	Execute(context.Context, types.Request, map[string]string) (string, error)
}

type executor struct {
	httpClient *http.Client
}

func NewExecutor(opts ...Option) Executor {
	e := &executor{
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func (e *executor) Execute(ctx context.Context, request types.Request, argValues map[string]string) (string, error) {
	endpoint, err := e.buildEndpoint(request, argValues)
	if err != nil {
		return "", err
	}

	fullURL := e.buildFullURL(request.Secure, request.Host, endpoint)
	body := e.buildRequestBody(argValues[request.Body])

	req, err := http.NewRequestWithContext(ctx, request.Method, fullURL, body)
	if err != nil {
		return "", err
	}

	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	slog.Info("executing http request",
		slog.Group("request",
			slog.String("method", request.Method),
			slog.String("url", fullURL),
		),
	)

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			slog.Error("failed to close response body")
		}
	}()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		slog.Error("http request failed",
			slog.Group("request",
				slog.String("method", request.Method),
				slog.String("url", fullURL),
				slog.String("response", string(bodyBytes)),
			),
		)
		return "", fmt.Errorf("http request failed: %s", bodyBytes)
	}

	response := types.Response{
		StatusCode: resp.StatusCode,
		Body:       string(bodyBytes),
	}

	result, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(result), nil
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(c *executor) {
		c.httpClient = httpClient
	}
}

func (e *executor) buildEndpoint(request types.Request, args map[string]string) (string, error) {
	endpoint := request.Endpoint

	for _, param := range request.PathParams {
		val, ok := args[param]
		if !ok {
			return "", fmt.Errorf("missing required path param: %s", param)
		}
		endpoint = strings.ReplaceAll(endpoint, ":"+param, url.PathEscape(val))
	}

	query := url.Values{}
	for _, key := range request.QueryParams {
		if val, ok := args[key]; ok {
			query.Set(key, val)
		}
	}

	if encoded := query.Encode(); encoded != "" {
		sep := "?"
		if strings.Contains(endpoint, "?") {
			sep = "&"
		}
		endpoint += sep + encoded
	}

	return endpoint, nil
}

func (e *executor) buildFullURL(secure bool, host, endpoint string) string {
	scheme := "http"
	if secure {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s%s", scheme, host, endpoint)
}

func (e *executor) buildRequestBody(body string) io.Reader {
	if body == "" {
		return nil
	}
	return strings.NewReader(body)
}
