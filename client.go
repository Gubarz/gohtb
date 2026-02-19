// Package gohtb provides the primary client for interacting with the Hack The Box API.
package gohtb

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	v4client "github.com/gubarz/gohtb/httpclient/v4"
	v5client "github.com/gubarz/gohtb/httpclient/v5"
	"github.com/gubarz/gohtb/internal/logging"
	"github.com/gubarz/gohtb/services/challenges"
	"github.com/gubarz/gohtb/services/containers"
	"github.com/gubarz/gohtb/services/fortresses"
	"github.com/gubarz/gohtb/services/machines"
	"github.com/gubarz/gohtb/services/prolabs"
	"github.com/gubarz/gohtb/services/rankings"
	"github.com/gubarz/gohtb/services/search"
	"github.com/gubarz/gohtb/services/seasons"
	"github.com/gubarz/gohtb/services/sherlocks"
	"github.com/gubarz/gohtb/services/teams"
	"github.com/gubarz/gohtb/services/users"
	"github.com/gubarz/gohtb/services/vms"
	"github.com/gubarz/gohtb/services/vpn"
)

// Client is the main API client for interacting with Hack The Box services.
// It holds configuration settings and provides access to various API endpoints
// through its service fields (e.g., Challengs, Machines, Seasons).
type Client struct {
	v4api       v4client.ClientInterface
	v5api       v5client.ClientInterface
	httpClient  *http.Client
	htbToken    string
	logger      Logger
	rateLimiter *RateLimiter
	server      string
	userAgent   string
	timeout     time.Duration
	debug       bool
	retryConfig RetryConfig

	// Services

	Challenges *challenges.Service
	Containers *containers.Service
	Fortresses *fortresses.Service
	Machines   *machines.Service
	Rankings   *rankings.Service
	Prolabs    *prolabs.Service
	Search     *search.Service
	Seasons    *seasons.Service
	Sherlocks  *sherlocks.Service
	Teams      *teams.Service
	Users      *users.Service
	// VMs is a service for managing virtual machines.
	// Can be used to Spawn, Stop, Extend, and Terminate VMs.
	VMs *vms.Service
	// VPN is a service for managing VPN connections and configurations.
	// This contains the endpoints for Access and Connections.
	VPN *vpn.Service
}

// Logger defines the logging interface used by the client.
// It's an alias for the internal logging.Logger interface.
type Logger = logging.Logger

// NoopLogger provides a Logger implementation that performs no operations.
// It's an alias for the internal logging.NoopLogger struct.
type NoopLogger = logging.NoopLogger

// RetryConfig specifies the configuration for automatic request retries.
type RetryConfig struct {
	// MaxRetries is the maximum number of retry attempts for a failed request.
	MaxRetries int
	// RetryPolicy determines whether a request should be retried and the wait duration.
	// If nil, DefaultRetryPolicy is used.
	RetryPolicy RetryPolicy
}

// Option defines the functional option type for configuring the Client.
type Option func(*Client)

const (
	baseHTBServer    = "https://labs.hackthebox.com/api"
	defaultUserAgent = "htb-go/" + version
	version          = "0.1"
)

// New creates and configures a new Hack The Box API Client.
// It requires a valid API token. Various aspects of the client can be configured
// by passing functional options (e.g., WithServer, WithLogger, WithTimeout).
//
// Example:
//
//	client, err := gohtb.New("YOUR_API_TOKEN",
//		gohtb.WithLogger(myCustomLogger),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Use the client...
//	info, err := client.Users.Info(context.Background())
func New(token string, options ...Option) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("htb token is required")
	}

	if err := isLikelyJWT(token); err != nil {
		return nil, err
	}

	c := &Client{
		htbToken:  token,
		server:    baseHTBServer,
		logger:    logging.NoopLogger{},
		userAgent: defaultUserAgent,
		timeout:   60 * time.Second,
		retryConfig: RetryConfig{
			MaxRetries:  4,
			RetryPolicy: &DefaultRetryPolicy{},
		},
	}

	for _, option := range options {
		option(c)
	}

	var finalHTTPClient *http.Client
	if c.httpClient != nil {
		finalHTTPClient = c.httpClient
		c.logger.Info("Using custom HTTP client provided via WithHTTPClient option. Note: Internal rate limiting and retry logic might be bypassed unless the custom client's transport is configured accordingly.")
		c.rateLimiter = NewRateLimiter(context.Background(), c.logger)

	} else {
		c.logger.Debug("Setting up default internal HTTP client with rate limiting and retries.")
		c.rateLimiter = NewRateLimiter(context.Background(), c.logger)
		apiTransport := NewAPITransport(
			http.DefaultTransport,
			c.rateLimiter,
			c.retryConfig,
			c.logger,
		)

		finalHTTPClient = &http.Client{
			Timeout:   c.timeout,
			Transport: apiTransport,
		}
		c.httpClient = finalHTTPClient
	}

	v4Server := c.server + "/v4"
	v4, err := v4client.NewClient(
		v4Server,
		v4client.WithHTTPClient(finalHTTPClient),
		v4client.WithRequestEditorFn(c.addHeaders),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	v5Server := c.server + "/v5"
	v5, err := v5client.NewClientWithResponses(
		v5Server,
		v5client.WithHTTPClient(finalHTTPClient),
		v5client.WithRequestEditorFn(c.addHeaders),
	)
	if err != nil {
		return nil, fmt.Errorf("init v5 client: %w", err)
	}

	c.v4api = v4
	c.v5api = v5
	wireServices(c)
	return c, nil
}

func (c *Client) addHeaders(ctx context.Context, req *http.Request) error {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.htbToken))
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	return nil
}

func wireServices(c *Client) {
	c.Challenges = challenges.NewService(c.asServiceClient(), "challenge")
	c.Containers = containers.NewService(c.asServiceClient())
	c.Fortresses = fortresses.NewService(c.asServiceClient())
	c.Machines = machines.NewService(c.asServiceClient(), "machine")
	c.Rankings = rankings.NewService(c.asServiceClient())
	c.Prolabs = prolabs.NewService(c.asServiceClient())
	c.Search = search.NewService(c.asServiceClient())
	c.Seasons = seasons.NewService(c.asServiceClient())
	c.Sherlocks = sherlocks.NewService(c.asServiceClient())
	c.Teams = teams.NewService(c.asServiceClient())
	c.Users = users.NewService(c.asServiceClient())
	c.VMs = vms.NewService(c.asServiceClient())
	c.VPN = vpn.NewService(c.asServiceClient())
}

// WithDebug enables or disables debug logging within the client's internal operations.
func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.debug = debug
	}
}

// WithLogger provides a custom logger implementation conforming to the Logger interface.
// By default, NoopLogger is used.
func WithLogger(logger Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithTimeout sets the request timeout for the internal HTTP client.
// Default is 60 seconds.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithRetry configures the automatic retry mechanism for requests.
func WithRetry(config RetryConfig) Option {
	return func(c *Client) {
		c.retryConfig = config
	}
}

// WithServer specifies a custom base URL for the Hack The Box API.
// Defaults to "https://labs.hackthebox.com/api".
// Do not include a trailing slash. v4 and v5 endpoints are derived from this base URL
func WithServer(server string) Option {
	return func(c *Client) {
		c.server = strings.TrimRight(server, "/")
	}
}

// WithUserAgent sets a custom User-Agent header for outgoing requests.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithHTTPClient allows providing a custom *http.Client.
// If provided, options like WithTimeout and the default transport setup
// (including rate limiting and retries via APITransport) will be bypassed.
// The provided client is used directly. The user is responsible for its configuration.
func WithHTTPClient(customClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = customClient
	}
}

// ExperimentalClient provides direct access to the generated OpenAPI clients.
//
// This is intended as an advanced escape hatch for unsupported endpoints or
// custom request/response handling.
type ExperimentalClient struct {
	client *Client
}

// V4 returns the generated v4 OpenAPI client.
func (e ExperimentalClient) V4() v4client.ClientInterface {
	return e.client.v4api
}

// V5 returns the generated v5 OpenAPI client.
func (e ExperimentalClient) V5() v5client.ClientInterface {
	return e.client.v5api
}

// WrapContext returns the context wrapped with the same limiter context used
// by service calls in this SDK.
func (e ExperimentalClient) WrapContext(ctx context.Context) context.Context {
	if e.client == nil || e.client.rateLimiter == nil {
		return ctx
	}
	return e.client.rateLimiter.Wrap(ctx)
}

// Experimental returns direct access to the underlying OpenAPI clients.
//
// WARNING: This is an advanced escape hatch for power users.
// The returned clients are auto-generated.
//
// This method is:
// - Unstable across versions (breaking changes will not be versioned or warned about)
// - Subject to change or removal at any time
// - Not covered by documentation or support
//
// Use at your own risk. If it breaks, you get to keep both pieces.
func (c *Client) Experimental() ExperimentalClient {
	return ExperimentalClient{client: c}
}

func isLikelyJWT(s string) error {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return errors.New("invalid token")
	}

	decodePart := func(part string) ([]byte, error) {
		decoded, err := base64.RawURLEncoding.DecodeString(part)
		if err != nil {
			return nil, errors.New("invalid token")
		}
		return decoded, nil
	}

	header, err := decodePart(parts[0])
	if err != nil {
		return errors.New("invalid token")
	}
	payload, err := decodePart(parts[1])
	if err != nil {
		return errors.New("invalid token")
	}

	if len(parts[2]) > 0 {
		if _, err := decodePart(parts[2]); err != nil {
			return errors.New("invalid token")
		}
	}

	if !json.Valid(header) || !json.Valid(payload) {
		return errors.New("invalid token")
	}

	return nil
}
