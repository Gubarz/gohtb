package seasons

import (
	"context"

	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	"github.com/gubarz/gohtb/internal/service"
)

func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

func (s *Service) Season(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

func (h *Handle) Rewards(ctx context.Context) (RewardsResponse, error) {
	resp, err := h.client.V4().GetSeasonRewardsWithResponse(h.client.Limiter().Wrap(ctx), h.id)
	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) RewardsResponse {
			return RewardsResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return RewardsResponse{
		Data: convert.SlicePointer(resp.JSON200.Data, fromAPISeasonRewardsDataItem),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) UserRank(ctx context.Context) (UserRankResponse, error) {
	resp, err := h.client.V4().GetSeasonUserRankWithResponse(h.client.Limiter().Wrap(ctx), h.id)
	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) UserRankResponse {
			return UserRankResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return UserRankResponse{
		Data: fromAPISeasonUserRankData(resp.JSON200.Data),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) UserFollowers(ctx context.Context) (UserFollowersResponse, error) {
	resp, err := h.client.V4().GetSeasonUserFollowersWithResponse(h.client.Limiter().Wrap(ctx), h.id)
	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) UserFollowersResponse {
			return UserFollowersResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return UserFollowersResponse{
		Data: fromAPISeasonUserFollowerData(resp.JSON200.Data),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (s *Service) List(ctx context.Context) (ListResponse, error) {
	resp, err := s.base.Client.V4().GetSeasonListWithResponse(s.base.Client.Limiter().Wrap(ctx))
	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) ListResponse {
			return ListResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return ListResponse{
		Data: convert.SlicePointer(resp.JSON200.Data, fromAPISeasonListDataItem),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (s *Service) Machines(ctx context.Context) (MachinesResponse, error) {
	resp, err := s.base.Client.V4().GetSeasonMachinesWithResponse(s.base.Client.Limiter().Wrap(ctx))
	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) MachinesResponse {
			return MachinesResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return MachinesResponse{
		Data: convert.SlicePointer(resp.JSON200.Data, fromAPISeasonMachinesDataItem),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (s *Service) ActiveMachine(ctx context.Context) (ActiveMachineResponse, error) {
	resp, err := s.base.Client.V4().GetSeasonMachineActiveWithResponse(s.base.Client.Limiter().Wrap(ctx))
	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) ActiveMachineResponse {
			return ActiveMachineResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return ActiveMachineResponse{
		Data: fromAPISeasonMachineActive(resp.JSON200.Data),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}
