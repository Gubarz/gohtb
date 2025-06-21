package testutil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	v4client "github.com/gubarz/gohtb/internal/httpclient/v4"
	"github.com/gubarz/gohtb/internal/logging"
	"github.com/stretchr/testify/require"
)

type TestLimiter struct{}

func (l *TestLimiter) Wrap(ctx context.Context) context.Context { return ctx }

type TestServiceClient struct {
	HttpClientInstance v4client.ClientWithResponsesInterface
}

func (t *TestServiceClient) HttpClient() v4client.ClientWithResponsesInterface {
	return t.HttpClientInstance
}
func (t *TestServiceClient) Limiter() interface {
	Wrap(context.Context) context.Context
} {
	return &TestLimiter{}
}
func (t *TestServiceClient) Logger() logging.Logger {
	// return logging.NoopLogger{}
	return nil
}

func NewTestServerAndClient(t *testing.T, handler http.HandlerFunc, opts ...v4client.ClientOption) (*httptest.Server, v4client.ClientWithResponsesInterface, func()) {
	ts := httptest.NewServer(handler)

	finalOpts := []v4client.ClientOption{v4client.WithHTTPClient(&http.Client{})}
	finalOpts = append(finalOpts, opts...)

	client, err := v4client.NewClientWithResponses(ts.URL, finalOpts...)
	require.NoError(t, err)

	cleanup := func() {
		ts.Close()
	}
	return ts, client, cleanup
}

func WithAcceptJSONRequestEditor() v4client.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Accept", "application/json")
		return nil
	}
}

func NewJSONSuccessHandler(t *testing.T, expectedPath string, responseBody []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
	}
}

func NewJSONErrorHandler(t *testing.T, expectedPath string, statusCode int, responseBody []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(responseBody)
	}
}
