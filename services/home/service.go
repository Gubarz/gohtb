package home

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

// Service provides access to home dashboard endpoints.
type Service struct {
	base service.Base
}

// NewService creates a new home service bound to a shared client.
//
// Example:
//
//	homeService := home.NewService(client)
//	_ = homeService
func NewService(client service.Client) *Service {
	return &Service{base: service.NewBase(client)}
}

type BannerData = v4Client.HomeBannersResponse

// BannerResponse contains home banner payload.
type BannerResponse struct {
	Data         BannerData
	ResponseMeta common.ResponseMeta
}

// Banner retrieves home banner content.
//
// Example:
//
//	banner, err := client.Home.Banner(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Home banner payload: %+v\n", banner.Data)
func (s *Service) Banner(ctx context.Context) (BannerResponse, error) {
	resp, err := s.base.Client.V4().GetHomeBanner(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return BannerResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetHomeBannerResponse)
	if err != nil {
		return BannerResponse{ResponseMeta: meta}, err
	}

	return BannerResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}
