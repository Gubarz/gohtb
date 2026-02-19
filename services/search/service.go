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

// NewService creates a new search service bound to a shared client.
//
// Example:
//
//	searchService := search.NewService(client)
//	_ = searchService
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

// Query creates a search handle for the provided keyword.
//
// Example:
//
//	query := client.Search.Query("kernel")
//	_ = query
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

// Users searches only user results for the configured query.
//
// Example:
//
//	users, err := client.Search.Query("john").Users(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("User hits: %d\n", len(users.Data.Users))
func (h *Handle) Users(ctx context.Context) (SearchResponse, error) {
	h.setTag(`["users"]`)
	return fetchSearch(ctx, h)
}

// Machines searches only machine results for the configured query.
//
// Example:
//
//	machines, err := client.Search.Query("forest").Machines(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Machine hits: %d\n", len(machines.Data.Machines))
func (h *Handle) Machines(ctx context.Context) (SearchResponse, error) {
	h.setTag(`["machines"]`)
	return fetchSearch(ctx, h)
}

// Challenges searches only challenge results for the configured query.
//
// Example:
//
//	challenges, err := client.Search.Query("crypto").Challenges(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Challenge hits: %d\n", len(challenges.Data.Challenges))
func (h *Handle) Challenges(ctx context.Context) (SearchResponse, error) {
	h.setTag(`["challenges"]`)
	return fetchSearch(ctx, h)
}

// Teams searches only team results for the configured query.
//
// Example:
//
//	teams, err := client.Search.Query("red").Teams(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Team hits: %d\n", len(teams.Data.Teams))
func (h *Handle) Teams(ctx context.Context) (SearchResponse, error) {
	h.setTag(`["teams"]`)
	return fetchSearch(ctx, h)
}

// All searches across all supported result types for the configured query.
//
// Example:
//
//	results, err := client.Search.Query("academy").All(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Search status: %t\n", results.Data.Status)
func (h *Handle) All(ctx context.Context) (SearchResponse, error) {
	h.tag = nil
	return fetchSearch(ctx, h)
}
