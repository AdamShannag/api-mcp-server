package request

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AdamShannag/api-mcp-server/pkg/types"
	"io"
	"log"
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
	endpoint := request.Endpoint

	for _, param := range request.PathParams {
		if val, ok := argValues[param]; ok {
			endpoint = strings.ReplaceAll(endpoint, ":"+param, url.PathEscape(val))
		} else {
			return "", fmt.Errorf("missing required path param: %s", param)
		}
	}

	query := url.Values{}
	for _, q := range request.QueryParams {
		if val, ok := argValues[q]; ok {
			query.Set(q, val)
		}
	}
	if encoded := query.Encode(); encoded != "" {
		endpoint += "?" + encoded
	}

	scheme := "http"
	if request.Secure {
		scheme = "https"
	}
	fullURL := fmt.Sprintf("%s://%s%s", scheme, request.Host, endpoint)

	var requestBody io.Reader
	if argValues[request.Body] != "" {
		requestBody = strings.NewReader(argValues[request.Body])
	}

	req, err := http.NewRequestWithContext(ctx, request.Method, fullURL, requestBody)
	if err != nil {
		return "", err
	}

	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	resp, err := e.httpClient.Do(req)

	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println("Failed to close response body")
		}
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	response := types.Response{
		StatusCode: resp.StatusCode,
		Body:       string(responseBody),
	}

	data, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(data), nil
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(c *executor) {
		c.httpClient = httpClient
	}
}
