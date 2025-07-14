package request_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/AdamShannag/api-mcp-server/pkg/request"
	"github.com/AdamShannag/api-mcp-server/pkg/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecute_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/todos/42", r.URL.Path)
		assert.Equal(t, "test", r.URL.Query().Get("q"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":42,"title":"Test Todo"}`))
	}))
	defer ts.Close()

	host := ts.URL[len("http://"):]

	req := types.Request{
		Method:      http.MethodGet,
		Host:        host,
		Endpoint:    "/todos/:id",
		Secure:      false,
		PathParams:  []string{"id"},
		QueryParams: []string{"q"},
	}

	ex := request.NewExecutor()
	result, err := ex.Execute(context.Background(), req, map[string]string{
		"id": "42",
		"q":  "test",
	})

	assert.NoError(t, err)

	var resp types.Response
	assert.NoError(t, json.Unmarshal([]byte(result), &resp))
	assert.Equal(t, 200, resp.StatusCode)
	assert.Contains(t, resp.Body, `"title":"Test Todo"`)
}

func TestExecute_MissingPathParam(t *testing.T) {
	req := types.Request{
		Method:     http.MethodGet,
		Host:       "example.com",
		Endpoint:   "/todos/:id",
		PathParams: []string{"id"},
	}

	ex := request.NewExecutor()
	_, err := ex.Execute(context.Background(), req, map[string]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required path param")
}

func TestExecute_ReadBodyFailure(t *testing.T) {
	req := types.Request{
		Method:   http.MethodGet,
		Host:     "api.test",
		Endpoint: "/fail",
		Secure:   false,
	}

	client := &http.Client{
		Transport: errorTransportFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       &errReader{err: errors.New("read error")},
				Header:     make(http.Header),
			}, nil
		}),
	}

	ex := request.NewExecutor(request.WithHttpClient(client))

	_, err := ex.Execute(context.Background(), req, map[string]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read response body")
}

type errReader struct {
	err error
}

func (e *errReader) Read(_ []byte) (int, error) { return 0, e.err }
func (e *errReader) Close() error               { return nil }

type errorTransportFunc func(*http.Request) (*http.Response, error)

func (f errorTransportFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
