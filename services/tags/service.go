package tags

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

// Service provides access to tag catalog endpoints.
type Service struct {
	base service.Base
}

// NewService creates a new tags service bound to a shared client.
//
// Example:
//
//	tagsService := tags.NewService(client)
//	_ = tagsService
func NewService(client service.Client) *Service {
	return &Service{base: service.NewBase(client)}
}

type ListData = v4Client.TagsListResponse

// ListResponse contains machine and challenge tag categories.
type ListResponse struct {
	Data         ListData
	ResponseMeta common.ResponseMeta
}

// List retrieves available tag categories and tags.
//
// Example:
//
//	tags, err := client.Tags.List(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Tags payload: %+v\n", tags.Data)
func (s *Service) List(ctx context.Context) (ListResponse, error) {
	resp, err := s.base.Client.V4().GetTagsList(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ListResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetTagsListResponse)
	if err != nil {
		return ListResponse{ResponseMeta: meta}, err
	}

	return ListResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}
