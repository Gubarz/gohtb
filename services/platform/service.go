package platform

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

// Service provides access to general platform endpoints.
type Service struct {
	base service.Base
}

// NewService creates a new platform service bound to a shared client.
//
// Example:
//
//	platformService := platform.NewService(client)
//	_ = platformService
func NewService(client service.Client) *Service {
	return &Service{base: service.NewBase(client)}
}

type NavigationMainData = v4Client.NavigationMainResponse

// NavigationMainResponse contains navigation summary data.
type NavigationMainResponse struct {
	Data         NavigationMainData
	ResponseMeta common.ResponseMeta
}

// NavigationMain retrieves main navigation summary information.
//
// Example:
//
//	navigation, err := client.Platform.NavigationMain(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Navigation payload: %+v\n", navigation.Data)
func (s *Service) NavigationMain(ctx context.Context) (NavigationMainResponse, error) {
	resp, err := s.base.Client.V4().GetNavigationMain(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return NavigationMainResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetNavigationMainResponse)
	if err != nil {
		return NavigationMainResponse{ResponseMeta: meta}, err
	}

	return NavigationMainResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type NoticesData = v4Client.NoticesResponse

// NoticesResponse contains active platform notices.
type NoticesResponse struct {
	Data         NoticesData
	ResponseMeta common.ResponseMeta
}

// Notices retrieves active platform notices.
//
// Example:
//
//	notices, err := client.Platform.Notices(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Notices payload: %+v\n", notices.Data)
func (s *Service) Notices(ctx context.Context) (NoticesResponse, error) {
	resp, err := s.base.Client.V4().GetNotices(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return NoticesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetNoticesResponse)
	if err != nil {
		return NoticesResponse{ResponseMeta: meta}, err
	}

	return NoticesResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}
