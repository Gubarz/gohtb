package search

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

type Service struct {
	base service.Base
}

func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

type Handle struct {
	client  service.Client
	keyword string
	tag     v4Client.SearchTags
}

func (s *Service) Query(keyword string) *Handle {
	return &Handle{
		client:  s.base.Client,
		keyword: keyword,
	}
}

type SearchFetchResponse = v4Client.SearchFetchResponse
type SearchResponse struct {
	Data         SearchFetchResponse
	ResponseMeta common.ResponseMeta
}

func fetchSearch(ctx context.Context, h *Handle) (SearchResponse, error) {
	params := &v4Client.GetSearchFetchParams{
		Query: h.keyword,
	}
	if h.tag != nil {
		params.Tags = &h.tag
	}
	resp, err := h.client.V4().GetSearchFetch(h.client.Limiter().Wrap(ctx), params)
	if err != nil {
		return SearchResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSearchFetchResponse)
	if err != nil {
		return SearchResponse{ResponseMeta: meta}, err
	}

	return SearchResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

func (h *Handle) setTag(json string) {
	h.tag = v4Client.SearchTags([]string{json})
}

func (h *Handle) Users(ctx context.Context) (SearchResponse, error) {
	h.setTag(`["users"]`)
	return fetchSearch(ctx, h)
}

func (h *Handle) Machines(ctx context.Context) (SearchResponse, error) {
	h.setTag(`["machines"]`)
	return fetchSearch(ctx, h)
}

func (h *Handle) Challenges(ctx context.Context) (SearchResponse, error) {
	h.setTag(`["challenges"]`)
	return fetchSearch(ctx, h)
}

func (h *Handle) Teams(ctx context.Context) (SearchResponse, error) {
	h.setTag(`["teams"]`)
	return fetchSearch(ctx, h)
}

func (h *Handle) All(ctx context.Context) (SearchResponse, error) {
	h.tag = nil
	return fetchSearch(ctx, h)
}
