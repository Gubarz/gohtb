package leaderboards

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/service"
)



func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

func (s *Service) Rankings() *Handle {
	return &Handle{
		client: s.base.Client,
	}
}

func (h *Handle) Users(ctx context.Context) (UserRankingsResponse, error) {
	// for `rankings/users` endpoint
	resp, err := h.client.V4().GetRankingsUsers(h.client.Limiter().Wrap(ctx))
	if err != nil {
		return UserRankingsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsUsersResponse)
	if err != nil {
		return UserRankingsResponse{ResponseMeta: meta}, err
	}
	return UserRankingsResponse{
		Data: convert.SlicePointer(parsed.JSON200.Data, fromLeaderBoardDataUsers), 
		ResponseMeta: meta,
	}, nil
}

func (h *Handle) Teams(ctx context.Context) (TeamRankingsResponse, error) {
	// for `rankings/teams` endpoint
	resp, err := h.client.V4().GetRankingsTeams(h.client.Limiter().Wrap(ctx))
	if err != nil {
		return TeamRankingsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsTeamsResponse)
	if err != nil {
		return TeamRankingsResponse{ResponseMeta: meta}, err
	}
	return TeamRankingsResponse{
		Data: convert.SlicePointer(parsed.JSON200.Data, fromLeaderBoardDataTeams), 
		ResponseMeta: meta,
	}, nil
}



func (h *Handle) Countries(ctx context.Context) (CountryRankingsResponse, error) {
	// for `rankings/countries` endpoint
	resp, err := h.client.V4().GetRankingsTeams(h.client.Limiter().Wrap(ctx))
	if err != nil {
		return CountryRankingsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsCountriesResponse)
	if err != nil {
		return CountryRankingsResponse{ResponseMeta: meta}, err
	}
	return CountryRankingsResponse{
		Data: convert.SlicePointer(parsed.JSON200.Data, fromLeaderBoardDataCountries), 
		ResponseMeta: meta,
	}, nil
}