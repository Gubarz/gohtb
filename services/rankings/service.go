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
