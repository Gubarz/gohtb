package subscriptions

import (
	"context"

	"github.com/gubarz/gohtb/httpclient/v5"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

// Service provides access to subscription products and prices.
type Service struct {
	base service.Base
}

// NewService returns a service bound to the API client.
func NewService(client service.Client) *Service {
	return &Service{base: service.NewBase(client)}
}

// PricingResponse contains subscription products, plans, prices, and HTTP response metadata.
type PricingResponse struct {
	Data         v5.SubscriptionPricingSourceResponse
	ResponseMeta common.ResponseMeta
}

// Pricing returns available subscription products, plans, and prices.
//
// Example:
//
//	pricing, err := client.Subscriptions.Pricing(ctx)
func (s *Service) Pricing(ctx context.Context) (PricingResponse, error) {
	client := s.base.Client
	resp, err := client.V5().GetSubscriptionPricing(client.Limiter().Wrap(ctx))
	if err != nil {
		return PricingResponse{}, err
	}
	parsed, meta, err := common.Parse(resp, v5.ParseGetSubscriptionPricingResponse)
	if err != nil {
		return PricingResponse{ResponseMeta: meta}, err
	}
	return PricingResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}
