package service

import (
	"context"

	v4client "github.com/gubarz/gohtb/httpclient/v4"
	v5client "github.com/gubarz/gohtb/httpclient/v5"
	"github.com/gubarz/gohtb/internal/logging"
)

// Client defines the common interface that all services expect
type Client interface {
	V4() v4client.ClientInterface
	V5() v5client.ClientInterface
	Limiter() interface {
		Wrap(context.Context) context.Context
	}
	Logger() logging.Logger
}

// Base provides common functionality for all services
type Base struct {
	Client Client
}

// NewBase creates a new base service
func NewBase(client Client) Base {
	return Base{Client: client}
}
