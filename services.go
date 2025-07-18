package gohtb

import (
	"context"

	v4client "github.com/gubarz/gohtb/httpclient/v4"
	v5client "github.com/gubarz/gohtb/httpclient/v5"
	"github.com/gubarz/gohtb/internal/logging"
)

type serviceAdapter struct {
	client *Client
}

// Create an adapter from the client
func (c *Client) asServiceClient() *serviceAdapter {
	return &serviceAdapter{client: c}
}

// Implement the service interface methods
func (a *serviceAdapter) V4() v4client.ClientInterface {
	return a.client.v4api
}

// Implement the service interface methods
func (a *serviceAdapter) V5() v5client.ClientInterface {
	return a.client.v5api
}

func (a *serviceAdapter) Limiter() interface {
	Wrap(context.Context) context.Context
} {
	return a.client.rateLimiter
}

func (a *serviceAdapter) Logger() logging.Logger {
	return a.client.logger
}
