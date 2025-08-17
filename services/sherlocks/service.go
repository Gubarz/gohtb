package sherlocks

import (
	"context"
	"strconv"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

func (s *Service) Sherlock(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

func (h *Handle) Info(ctx context.Context) (InfoResponse, error) {
	slug := strconv.Itoa(h.id)
	resp, err := h.client.V4().GetSherlock(h.client.Limiter().Wrap(ctx), slug)

	if err != nil {
		return InfoResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSherlockResponse)
	if err != nil {
		return InfoResponse{ResponseMeta: meta}, err
	}

	return InfoResponse{
		Data:         fromAPISherlock(parsed.JSON200.Data),
		ResponseMeta: meta,
	}, nil
}
