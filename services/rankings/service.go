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

// NewService creates a new rankings service bound to a shared client.
//
// Example:
//
//	rankingsService := rankings.NewService(client)
//	_ = rankingsService
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

// Countries retrieves country rankings.
//
// Example:
//
//	countries, err := client.Rankings.Countries(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Countries ranked: %d\n", len(countries.Data))
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

// Teams retrieves team rankings.
//
// Example:
//
//	teams, err := client.Rankings.Teams(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Teams ranked: %d\n", len(teams.Data))
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

// Users retrieves user rankings.
//
// Example:
//
//	users, err := client.Rankings.Users(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Users ranked: %d\n", len(users.Data))
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

// Country returns a handle for country-scoped ranking endpoints.
//
// Example:
//
//	country := client.Rankings.Country("US")
//	_ = country
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

// Members retrieves member ranking data for the selected country.
//
// Example:
//
//	members, err := client.Rankings.Country("US").Members(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Country members: %d\n", len(members.Data.Items))
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

type CurrentTeam struct {
	client service.Client
}

// CurrentTeam returns a handle for current-team ranking endpoints.
//
// Example:
//
//	team := client.Rankings.CurrentTeam()
//	_ = team
func (s *Service) CurrentTeam() *CurrentTeam {
	return &CurrentTeam{
		client: s.base.Client,
	}
}

type Team struct {
	client service.Client
	id     int
}

// Team returns a handle for team ranking endpoints scoped to a specific team ID.
//
// Example:
//
//	team := client.Rankings.Team(12345)
//	_ = team
func (s *Service) Team(id int) *Team {
	return &Team{
		client: s.base.Client,
		id:     id,
	}
}

// Best retrieves global current-team best-history ranking data.
//
// Example:
//
//	best, err := client.Rankings.CurrentTeam().Best(ctx, "1Y")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Team best rank: %d\n", best.Data.Rank)
func (t *CurrentTeam) Best(ctx context.Context, period string) (TeamBestResponse, error) {
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

// Overview retrieves global current-team trend metrics.
//
// Example:
//
//	overview, err := client.Rankings.CurrentTeam().Overview(ctx, "1Y")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Team points diff: %d\n", overview.Data.PointsDiff)
func (t *CurrentTeam) Overview(ctx context.Context, period string) (TeamOverviewResponse, error) {
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

// RankingBracket retrieves global current-team ranking bracket data.
//
// Example:
//
//	bracket, err := client.Rankings.CurrentTeam().RankingBracket(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Team bracket: %s\n", bracket.Data.CurrentBracket)
func (t *CurrentTeam) RankingBracket(ctx context.Context) (TeamBracketResponse, error) {
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

type RankingsTeamBestData = v4Client.RankingsTeamBestData
type TeamBestResponse struct {
	Data         RankingsTeamBestData
	ResponseMeta common.ResponseMeta
}

// Best retrieves best-history ranking data for the selected team.
//
// Example:
//
//	best, err := client.Rankings.Team(12345).Best(ctx, "1Y")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Team best rank: %d\n", best.Data.Rank)
func (t *Team) Best(ctx context.Context, period string) (TeamBestResponse, error) {
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

// Overview retrieves trend metrics for the selected team.
//
// Example:
//
//	overview, err := client.Rankings.Team(12345).Overview(ctx, "1Y")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Team points growth: %s\n", overview.Data.PointsGrowth)
func (t *Team) Overview(ctx context.Context, period string) (TeamOverviewResponse, error) {
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

// RankingBracket retrieves ranking bracket data for the selected team.
//
// Example:
//
//	bracket, err := client.Rankings.Team(12345).RankingBracket(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Team rank: %d\n", bracket.Data.Rank)
func (t *Team) RankingBracket(ctx context.Context) (TeamBracketResponse, error) {
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

type RankingsData = v4Client.RankingsData

type OverviewResponse struct {
	Data         RankingsData
	ResponseMeta common.ResponseMeta
}

// Overview retrieves the combined team and user ranking overview.
//
// Example:
//
//	overview, err := client.Rankings.Overview(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("User rank: %d\n", overview.Data.User.Rank)
func (s *Service) Overview(ctx context.Context) (OverviewResponse, error) {
	resp, err := s.base.Client.V4().GetRankings(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return OverviewResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsResponse)
	if err != nil {
		return OverviewResponse{ResponseMeta: meta}, err
	}

	return OverviewResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type RankingsCountryBestData = v4Client.RankingsCountryData

type CountryBestResponse struct {
	Data         RankingsCountryBestData
	ResponseMeta common.ResponseMeta
}

// CountryBest retrieves country ranking best-history data.
//
// Example:
//
//	best, err := client.Rankings.CountryBest(ctx, "1Y")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Country best rank: %d\n", best.Data.Rank)
func (s *Service) CountryBest(ctx context.Context, period string) (CountryBestResponse, error) {
	p := v4Client.GetRankingsCountryBestParamsPeriodN1Y
	if period != "" {
		p = v4Client.GetRankingsCountryBestParamsPeriod(period)
	}

	params := &v4Client.GetRankingsCountryBestParams{
		Period: p,
	}

	resp, err := s.base.Client.V4().GetRankingsCountryBest(s.base.Client.Limiter().Wrap(ctx), params)
	if err != nil {
		return CountryBestResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsCountryBestResponse)
	if err != nil {
		return CountryBestResponse{ResponseMeta: meta}, err
	}

	return CountryBestResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type RankingsCountryOverviewData = v4Client.RankingsCountryOverviewData

type CountryOverviewResponse struct {
	Data         RankingsCountryOverviewData
	ResponseMeta common.ResponseMeta
}

// CountryOverview retrieves country ranking trend metrics.
//
// Example:
//
//	overview, err := client.Rankings.CountryOverview(ctx, "1Y")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Country points diff: %d\n", overview.Data.PointsDiff)
func (s *Service) CountryOverview(ctx context.Context, period string) (CountryOverviewResponse, error) {
	p := v4Client.GetRankingsCountryOverviewParamsPeriodN1Y
	if period != "" {
		p = v4Client.GetRankingsCountryOverviewParamsPeriod(period)
	}

	params := &v4Client.GetRankingsCountryOverviewParams{
		Period: p,
	}

	resp, err := s.base.Client.V4().GetRankingsCountryOverview(s.base.Client.Limiter().Wrap(ctx), params)
	if err != nil {
		return CountryOverviewResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsCountryOverviewResponse)
	if err != nil {
		return CountryOverviewResponse{ResponseMeta: meta}, err
	}

	return CountryOverviewResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type RankingsCountryBracketData = v4Client.RankingsBracketData

type CountryBracketResponse struct {
	Data         RankingsCountryBracketData
	ResponseMeta common.ResponseMeta
}

// CountryRankingBracket retrieves bracket data for country rankings.
//
// Example:
//
//	bracket, err := client.Rankings.CountryRankingBracket(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Country bracket: %s\n", bracket.Data.CurrentBracket)
func (s *Service) CountryRankingBracket(ctx context.Context) (CountryBracketResponse, error) {
	resp, err := s.base.Client.V4().GetRankingsCountryRankingBracket(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return CountryBracketResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsCountryRankingBracketResponse)
	if err != nil {
		return CountryBracketResponse{ResponseMeta: meta}, err
	}

	return CountryBracketResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type RankingsUniversityItem = v4Client.RankingsUniversitiesItem
type RankingsUniversityItems = []RankingsUniversityItem

type UniversityRankingsResponse struct {
	Data         RankingsUniversityItems
	ResponseMeta common.ResponseMeta
}

// Universities retrieves university rankings.
//
// Example:
//
//	universities, err := client.Rankings.Universities(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Universities found: %d\n", len(universities.Data))
func (s *Service) Universities(ctx context.Context) (UniversityRankingsResponse, error) {
	resp, err := s.base.Client.V4().GetRankingsUniversities(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return UniversityRankingsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsUniversitiesResponse)
	if err != nil {
		return UniversityRankingsResponse{ResponseMeta: meta}, err
	}

	return UniversityRankingsResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type University struct {
	client service.Client
	id     int
}

// University returns a handle for ranking endpoints scoped to a university ID.
//
// Example:
//
//	university := client.Rankings.University(123)
//	_ = university
func (s *Service) University(id int) *University {
	return &University{
		client: s.base.Client,
		id:     id,
	}
}

type RankingsUniversityBracketData = v4Client.RankingsBracketData

type UniversityBracketResponse struct {
	Data         RankingsUniversityBracketData
	ResponseMeta common.ResponseMeta
}

// RankingBracket retrieves ranking bracket data for this university.
//
// Example:
//
//	bracket, err := client.Rankings.University(123).RankingBracket(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("University rank: %d\n", bracket.Data.Rank)
func (u *University) RankingBracket(ctx context.Context) (UniversityBracketResponse, error) {
	resp, err := u.client.V4().GetRankingsUniversityRankingBracketId(
		u.client.Limiter().Wrap(ctx),
		v4Client.UniversityId(u.id),
	)
	if err != nil {
		return UniversityBracketResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsUniversityRankingBracketIdResponse)
	if err != nil {
		return UniversityBracketResponse{ResponseMeta: meta}, err
	}

	return UniversityBracketResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type RankingsUserBestData = v4Client.RankingsUserBestData

type UserBestResponse struct {
	Data         RankingsUserBestData
	ResponseMeta common.ResponseMeta
}

// UserBest retrieves user ranking best-history data.
//
// Example:
//
//	best, err := client.Rankings.UserBest(ctx, "1Y")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Best user rank: %.0f\n", best.Data.Rank)
func (s *Service) UserBest(ctx context.Context, period string) (UserBestResponse, error) {
	p := v4Client.GetRankingsUserBestParamsPeriodN1Y
	if period != "" {
		p = v4Client.GetRankingsUserBestParamsPeriod(period)
	}

	params := &v4Client.GetRankingsUserBestParams{
		Period: p,
	}

	resp, err := s.base.Client.V4().GetRankingsUserBest(s.base.Client.Limiter().Wrap(ctx), params)
	if err != nil {
		return UserBestResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsUserBestResponse)
	if err != nil {
		return UserBestResponse{ResponseMeta: meta}, err
	}

	return UserBestResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type RankingsUserOverviewData = v4Client.RankingsUserOverviewData

type UserOverviewResponse struct {
	Data         RankingsUserOverviewData
	ResponseMeta common.ResponseMeta
}

// UserOverview retrieves user ranking trend metrics.
//
// Example:
//
//	overview, err := client.Rankings.UserOverview(ctx, "1Y")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("User points growth: %s\n", overview.Data.PointsGrowth)
func (s *Service) UserOverview(ctx context.Context, period string) (UserOverviewResponse, error) {
	p := v4Client.GetRankingsUserOverviewParamsPeriodN1Y
	if period != "" {
		p = v4Client.GetRankingsUserOverviewParamsPeriod(period)
	}

	params := &v4Client.GetRankingsUserOverviewParams{
		Period: p,
	}

	resp, err := s.base.Client.V4().GetRankingsUserOverview(s.base.Client.Limiter().Wrap(ctx), params)
	if err != nil {
		return UserOverviewResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsUserOverviewResponse)
	if err != nil {
		return UserOverviewResponse{ResponseMeta: meta}, err
	}

	return UserOverviewResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type RankingsUserBracketData = v4Client.RankingUserRankBracketData

type UserBracketResponse struct {
	Data         RankingsUserBracketData
	ResponseMeta common.ResponseMeta
}

// UserRankingBracket retrieves user ranking bracket data.
//
// Example:
//
//	bracket, err := client.Rankings.UserRankingBracket(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("User bracket: %s\n", bracket.Data.CurrentBracket)
func (s *Service) UserRankingBracket(ctx context.Context) (UserBracketResponse, error) {
	resp, err := s.base.Client.V4().GetRankingsUserRankingBracket(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return UserBracketResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetRankingsUserRankingBracketResponse)
	if err != nil {
		return UserBracketResponse{ResponseMeta: meta}, err
	}

	return UserBracketResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}
