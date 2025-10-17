package rankings

import (
	"context"
	"errors"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/errutil"
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

type RankingsCountriesItem = v4Client.RankingsCountriesItem
type RankingsCountriesItems = []RankingsCountriesItem

type CountryRankingsResponse struct {
	Data         RankingsCountriesItems
	ResponseMeta common.ResponseMeta
}

func (s *Service) Countries(ctx context.Context) (CountryRankingsResponse, error) {
	resp, err := s.base.Client.V4().GetRankingsCountries(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return CountryRankingsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsCountriesResponse)
	if err != nil {
		return CountryRankingsResponse{ResponseMeta: meta}, err
	}
	return CountryRankingsResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type RankingTeamItem = v4Client.RankingsTeamItem
type RankingsTeamItems = []RankingTeamItem

type TeamRankingsResponse struct {
	Data         RankingsTeamItems
	ResponseMeta common.ResponseMeta
}

func (s *Service) Teams(ctx context.Context) (TeamRankingsResponse, error) {
	resp, err := s.base.Client.V4().GetRankingsTeams(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return TeamRankingsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsTeamsResponse)
	if err != nil {
		return TeamRankingsResponse{ResponseMeta: meta}, err
	}
	return TeamRankingsResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type RankingUserDataItem = v4Client.RankingsUserData
type RankingUserDataItems = []RankingUserDataItem

type UserRankingsResponse struct {
	Data         RankingUserDataItems
	ResponseMeta common.ResponseMeta
}

func (s *Service) Users(ctx context.Context) (UserRankingsResponse, error) {
	resp, err := s.base.Client.V4().GetRankingsUsers(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return UserRankingsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsUsersResponse)
	if err != nil {
		return UserRankingsResponse{ResponseMeta: meta}, err
	}
	return UserRankingsResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type Country struct {
	client    service.Client
	shortName string
}

func (s *Service) Country(shortName string) *Country {
	return &Country{
		client:    s.base.Client,
		shortName: shortName,
	}
}

type RankingsCountryByMembersData = v4Client.RankingCountryMemberData

type CountryRankingsByMembersResponse struct {
	Data         RankingsCountryByMembersData
	ResponseMeta common.ResponseMeta
}

func (c *Country) Members(ctx context.Context) (CountryRankingsByMembersResponse, error) {
	if c.shortName == "" {
		_, apiErr := errutil.UnwrapFailure[common.ResponseMeta](errors.New("country short name is required"), nil, 0, nil)
		return CountryRankingsByMembersResponse{ResponseMeta: common.ResponseMeta{}}, apiErr
	}
	resp, err := c.client.V4().GetRankingsCountryUSMembers(c.client.Limiter().Wrap(ctx),
		c.shortName)
	if err != nil {
		return CountryRankingsByMembersResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsCountryUSMembersResponse)
	if err != nil {
		return CountryRankingsByMembersResponse{ResponseMeta: meta}, err
	}
	return CountryRankingsByMembersResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type Team struct {
	client service.Client
}

func (s *Service) Team(id int) *Team {
	return &Team{
		client: s.base.Client,
	}
}

func (t *Team) Best(ctx context.Context, period string) (TeamBestResponse, error) {
	var p v4Client.GetRankingsTeamBestParamsPeriod = "1Y"
	if period != "" {
		p = v4Client.GetRankingsTeamBestParamsPeriod(period)
	}
	params := v4Client.GetRankingsTeamBestParams{
		Period: p,
	}
	resp, err := t.client.V4().GetRankingsTeamBest(t.client.Limiter().Wrap(ctx),
		&params)
	if err != nil {
		return TeamBestResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsTeamBestResponse)
	if err != nil {
		return TeamBestResponse{ResponseMeta: meta}, err
	}
	return TeamBestResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

func (t *Team) Overview(ctx context.Context, period string) (TeamOverviewResponse, error) {
	var p v4Client.GetRankingsTeamOverviewParamsPeriod = "1Y"
	if period != "" {
		p = v4Client.GetRankingsTeamOverviewParamsPeriod(period)
	}
	params := v4Client.GetRankingsTeamOverviewParams{
		Period: p,
	}
	resp, err := t.client.V4().GetRankingsTeamOverview(t.client.Limiter().Wrap(ctx),
		&params)
	if err != nil {
		return TeamOverviewResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsTeamOverviewResponse)
	if err != nil {
		return TeamOverviewResponse{ResponseMeta: meta}, err
	}
	return TeamOverviewResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

func (t *Team) RankingBracket(ctx context.Context) (TeamBracketResponse, error) {
	resp, err := t.client.V4().GetRankingsTeamRankingBracket(t.client.Limiter().Wrap(ctx))
	if err != nil {
		return TeamBracketResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsTeamRankingBracketResponse)
	if err != nil {
		return TeamBracketResponse{ResponseMeta: meta}, err
	}
	return TeamBracketResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

func (t *Team) ByID(id int) *TeamByID {
	return &TeamByID{
		client: t.client,
		id:     id,
	}
}

type TeamByID struct {
	client service.Client
	id     int
}

type RankingsTeamBestData = v4Client.RankingsTeamBestData
type TeamBestResponse struct {
	Data         RankingsTeamBestData
	ResponseMeta common.ResponseMeta
}

func (t *TeamByID) Best(ctx context.Context, period string) (TeamBestResponse, error) {
	var p v4Client.GetRankingsTeamBestIdParamsPeriod = "1Y"
	if period != "" {
		p = v4Client.GetRankingsTeamBestIdParamsPeriod(period)
	}
	params := v4Client.GetRankingsTeamBestIdParams{
		Period: p,
	}
	resp, err := t.client.V4().GetRankingsTeamBestId(t.client.Limiter().Wrap(ctx),
		v4Client.TeamId(t.id), &params)
	if err != nil {
		return TeamBestResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsTeamBestIdResponse)
	if err != nil {
		return TeamBestResponse{ResponseMeta: meta}, err
	}
	return TeamBestResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type RankingsTeamOverviewData = v4Client.RankingsTeamOverviewData

type TeamOverviewResponse struct {
	Data         RankingsTeamOverviewData
	ResponseMeta common.ResponseMeta
}

func (t *TeamByID) Overview(ctx context.Context, period string) (TeamOverviewResponse, error) {
	var p v4Client.GetRankingsTeamOverviewIdParamsPeriod = "1Y"
	if period != "" {
		p = v4Client.GetRankingsTeamOverviewIdParamsPeriod(period)
	}
	params := v4Client.GetRankingsTeamOverviewIdParams{
		Period: p,
	}
	resp, err := t.client.V4().GetRankingsTeamOverviewId(t.client.Limiter().Wrap(ctx),
		v4Client.TeamId(t.id), &params)
	if err != nil {
		return TeamOverviewResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsTeamOverviewIdResponse)
	if err != nil {
		return TeamOverviewResponse{ResponseMeta: meta}, err
	}
	return TeamOverviewResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type RankingsTeamRankingBracketData = v4Client.RankingsTeamRankingBracketData

type TeamBracketResponse struct {
	Data         RankingsTeamRankingBracketData
	ResponseMeta common.ResponseMeta
}

func (t *TeamByID) RankingBracket(ctx context.Context) (TeamBracketResponse, error) {
	resp, err := t.client.V4().GetRankingsTeamRankingBracketId(t.client.Limiter().Wrap(ctx),
		v4Client.TeamId(t.id))
	if err != nil {
		return TeamBracketResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsTeamRankingBracketIdResponse)
	if err != nil {
		return TeamBracketResponse{ResponseMeta: meta}, err
	}
	return TeamBracketResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}
