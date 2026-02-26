package pwnbox

import (
	"context"
	"fmt"
	"net/http"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	"github.com/gubarz/gohtb/internal/service"
)

// Service provides access to Pwnbox lifecycle endpoints.
type Service struct {
	base service.Base
}

// NewService creates a new pwnbox service bound to a shared client.
//
// Example:
//
//	pwnboxService := pwnbox.NewService(client)
//	_ = pwnboxService
func NewService(client service.Client) *Service {
	return &Service{base: service.NewBase(client)}
}

type StartData = v4Client.PwnboxStartResponse

// StartResponse contains pwnbox start result.
type StartResponse struct {
	Data         StartData
	ResponseMeta common.ResponseMeta
}

// Start starts a new Pwnbox instance.
//
// Example:
//
//	start, err := client.Pwnbox.Start(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Pwnbox start response: %+v\n", start.Data)
func (s *Service) Start(ctx context.Context) (StartResponse, error) {
	resp, err := s.base.Client.V4().PostPwnboxStart(
		s.base.Client.Limiter().Wrap(ctx),
		v4Client.PostPwnboxStartJSONRequestBody{},
	)
	if err != nil {
		return StartResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostPwnboxStartResponse)
	if err != nil {
		return StartResponse{ResponseMeta: meta}, err
	}

	return StartResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type StatusData = v4Client.PwnboxStatusResponse

// StatusResponse contains current Pwnbox state.
type StatusResponse struct {
	Data         StatusData
	ResponseMeta common.ResponseMeta
}

// Status retrieves current Pwnbox status.
//
// Example:
//
//	status, err := client.Pwnbox.Status(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Pwnbox status payload: %+v\n", status.Data)
func (s *Service) Status(ctx context.Context) (StatusResponse, error) {
	resp, err := s.base.Client.V4().GetPwnboxStatus(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return StatusResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetPwnboxStatusResponse)
	if err != nil {
		return StatusResponse{ResponseMeta: meta}, err
	}

	return StatusResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

// Terminate stops the active Pwnbox instance.
//
// Example:
//
//	result, err := client.Pwnbox.Terminate(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Terminate result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (s *Service) Terminate(ctx context.Context) (common.MessageResponse, error) {
	resp, err := s.base.Client.V4().PostPwnboxTerminate(s.base.Client.Limiter().Wrap(ctx))
	raw := extract.Raw(resp)
	if err != nil || resp == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) common.MessageResponse {
			return common.MessageResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return errutil.UnwrapFailure(fmt.Errorf("unexpected status code %d", resp.StatusCode), raw, resp.StatusCode, func(raw []byte) common.MessageResponse {
			return common.MessageResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return common.MessageResponse{
		Data: common.Message{Message: http.StatusText(resp.StatusCode), Success: true},
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			CFRay:      resp.Header.Get("CF-Ray"),
		},
	}, nil
}

type UsageData = v4Client.PwnboxUsageResponse

// UsageResponse contains pwnbox usage counters.
type UsageResponse struct {
	Data         UsageData
	ResponseMeta common.ResponseMeta
}

// Usage retrieves Pwnbox usage information.
//
// Example:
//
//	usage, err := client.Pwnbox.Usage(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Pwnbox usage payload: %+v\n", usage.Data)
func (s *Service) Usage(ctx context.Context) (UsageResponse, error) {
	resp, err := s.base.Client.V4().GetPwnboxUsage(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return UsageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetPwnboxUsageResponse)
	if err != nil {
		return UsageResponse{ResponseMeta: meta}, err
	}

	return UsageResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}
