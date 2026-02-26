package badges

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

// Service provides access to badge catalog endpoints.
type Service struct {
	base service.Base
}

// NewService creates a new badges service bound to a shared client.
//
// Example:
//
//	badgesService := badges.NewService(client)
//	_ = badgesService
func NewService(client service.Client) *Service {
	return &Service{base: service.NewBase(client)}
}

type ListData = v4Client.BadgesResponse

// ListResponse contains badge catalog data.
type ListResponse struct {
	Data         ListData
	ResponseMeta common.ResponseMeta
}

// List retrieves the badge catalog.
//
// Example:
//
//	badges, err := client.Badges.List(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Badges payload: %+v\n", badges.Data)
func (s *Service) List(ctx context.Context) (ListResponse, error) {
	resp, err := s.base.Client.V4().GetBadges(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ListResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetBadgesResponse)
	if err != nil {
		return ListResponse{ResponseMeta: meta}, err
	}

	return ListResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}
