package universities

import (
	"context"

	"github.com/gubarz/gohtb/httpclient/v5"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

// Service provides access to university listings and statistics.
type Service struct {
	base service.Base
}

// NewService returns a service bound to the API client.
func NewService(client service.Client) *Service {
	return &Service{base: service.NewBase(client)}
}

// Handle selects a university for further queries.
type Handle struct {
	client service.Client
	id     int
}

// University selects a university by ID.
func (s *Service) University(id int) *Handle {
	return &Handle{client: s.base.Client, id: id}
}

// ListResponse contains a page of universities and HTTP response metadata.
type ListResponse struct {
	Data         v5.UniversitiesResponse
	ResponseMeta common.ResponseMeta
}

// List returns a page of universities.
//
// Example:
//
//	page, err := client.Universities.List(ctx, nil)
func (s *Service) List(ctx context.Context, params *v5.GetUniversitiesParams) (ListResponse, error) {
	client := s.base.Client
	resp, err := client.V5().GetUniversities(client.Limiter().Wrap(ctx), params)
	if err != nil {
		return ListResponse{}, err
	}
	parsed, meta, err := common.Parse(resp, v5.ParseGetUniversitiesResponse)
	if err != nil {
		return ListResponse{ResponseMeta: meta}, err
	}
	return ListResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

// StatsResponse contains university statistics and HTTP response metadata.
type StatsResponse struct {
	Data         v5.UniversityStatsResponse
	ResponseMeta common.ResponseMeta
}

// Stats returns statistics for this university.
//
// Example:
//
//	stats, err := client.Universities.University(42).Stats(ctx)
func (h *Handle) Stats(ctx context.Context) (StatsResponse, error) {
	client := h.client
	resp, err := client.V5().GetUniversityStats(client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return StatsResponse{}, err
	}
	parsed, meta, err := common.Parse(resp, v5.ParseGetUniversityStatsResponse)
	if err != nil {
		return StatsResponse{ResponseMeta: meta}, err
	}
	return StatsResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}
