package seasons

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	v5Client "github.com/gubarz/gohtb/httpclient/v5"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

// Service provides access to seasonal content and rankings.
type Service struct {
	base service.Base
}

// NewService creates a new seasons service bound to a shared client.
//
// Example:
//
//	seasonService := seasons.NewService(client)
//	_ = seasonService
func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

// Handle selects a season for content, rankings, and rewards.
type Handle struct {
	client service.Client
	id     int
}

// Season returns a handle for a specific season with the given ID.
// This handle can be used to perform operations related to that season,
// such as retrieving rewards, user rankings, and follower information.
//
// Example:
//
//	season := client.Seasons.Season(3)
//	_ = season
func (s *Service) Season(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

// RewardsResponse contains the seasonal reward groups and HTTP response metadata.
type RewardsResponse struct {
	Data         []v5Client.SeasonRewardsResponse_Data_Item
	ResponseMeta common.ResponseMeta
}

// Rewards returns the rewards for this season.
//
// Example:
//
//	result, err := client.Seasons.Season(21).Rewards(ctx)
func (h *Handle) Rewards(ctx context.Context) (RewardsResponse, error) {
	client := h.client
	resp, err := client.V5().GetSeasonRewards(client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return RewardsResponse{}, err
	}
	parsed, meta, err := common.Parse(resp, v5Client.ParseGetSeasonRewardsResponse)
	if err != nil {
		return RewardsResponse{ResponseMeta: meta}, err
	}
	return RewardsResponse{Data: parsed.JSON200.Data, ResponseMeta: meta}, nil
}

// UserRankResponse contains the signed-in user’s league progress and HTTP response metadata.
type UserRankResponse struct {
	Data         v5Client.SeasonUserRankResponse_Data
	ResponseMeta common.ResponseMeta
}

// UserRank returns the signed-in user’s league rankings and progress.
//
// Example:
//
//	result, err := client.Seasons.Season(21).UserRank(ctx)
func (h *Handle) UserRank(ctx context.Context) (UserRankResponse, error) {
	client := h.client
	resp, err := client.V5().GetSeasonUserRank(client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return UserRankResponse{}, err
	}
	parsed, meta, err := common.Parse(resp, v5Client.ParseGetSeasonUserRankResponse)
	if err != nil {
		return UserRankResponse{ResponseMeta: meta}, err
	}
	return UserRankResponse{Data: parsed.JSON200.Data, ResponseMeta: meta}, nil
}

type SeasonListDataItem = v4Client.SeasonListDataItem

type ListResponse struct {
	Data         []SeasonListDataItem
	ResponseMeta common.ResponseMeta
}

// List retrieves all available seasons on the HackTheBox platform.
// This returns a comprehensive list of all seasons, including current and past seasons.
//
// Example:
//
//	seasons, err := client.Seasons.List(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, season := range seasons.Data {
//		fmt.Printf("Season: %s (ID: %d)\n", season.Name, season.Id)
//	}
func (s *Service) List(ctx context.Context) (ListResponse, error) {
	resp, err := s.base.Client.V4().GetSeasonList(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ListResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSeasonListResponse)
	if err != nil {
		return ListResponse{ResponseMeta: meta}, err
	}

	return ListResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type SeasonMachinesDataItem = v4Client.SeasonMachinesDataItem

type MachinesResponse struct {
	Data         []SeasonMachinesDataItem
	ResponseMeta common.ResponseMeta
}

// Machines retrieves all machines available in the current season.
// This returns information about machines that are part of the active season,
// including their difficulty, points, and availability status.
//
// Example:
//
//	machines, err := client.Seasons.Season(3).Machines(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, machine := range machines.Data {
//		fmt.Printf("Machine: %s (Difficulty: %s)\n", machine.Name, machine.DifficultyText)
//	}
func (h *Handle) Machines(ctx context.Context) (MachinesResponse, error) {
	resp, err := h.client.V4().GetSeasonMachines(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return MachinesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSeasonMachinesResponse)
	if err != nil {
		return MachinesResponse{ResponseMeta: meta}, err
	}

	return MachinesResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

// ContentResponse contains the seasonal content and HTTP response metadata.
type ContentResponse struct {
	Data         v5Client.SeasonContentResponse_Data
	ResponseMeta common.ResponseMeta
}

// Content returns the seasonal content for this season.
//
// Example:
//
//	result, err := client.Seasons.Season(21).Content(ctx)
func (h *Handle) Content(ctx context.Context) (ContentResponse, error) {
	client := h.client
	resp, err := client.V5().GetSeasonContent(client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return ContentResponse{}, err
	}
	parsed, meta, err := common.Parse(resp, v5Client.ParseGetSeasonContentResponse)
	if err != nil {
		return ContentResponse{ResponseMeta: meta}, err
	}
	return ContentResponse{Data: parsed.JSON200.Data, ResponseMeta: meta}, nil
}

type SeasonActiveData = v4Client.SeasonActiveData

type ActiveMachineResponse struct {
	Data         SeasonActiveData
	ResponseMeta common.ResponseMeta
}

// ActiveMachine retrieves information about the currently active machine in the season.
// This returns details about the machine that is currently available for solving
// in the active season.
//
// Example:
//
//	activeMachine, err := client.Seasons.ActiveMachine(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Active machine: %s (ID: %d)\n", activeMachine.Data.Name, activeMachine.Data.Id)
func (s *Service) ActiveMachine(ctx context.Context) (ActiveMachineResponse, error) {
	resp, err := s.base.Client.V4().GetSeasonMachineActive(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ActiveMachineResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSeasonMachineActiveResponse)
	if err != nil {
		return ActiveMachineResponse{ResponseMeta: meta}, err
	}

	return ActiveMachineResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

// ActiveContentResponse contains current seasonal content and HTTP response metadata.
type ActiveContentResponse struct {
	Data         v5Client.SeasonActiveContentResponse_Data
	ResponseMeta common.ResponseMeta
}

// ActiveContent returns current offensive and defensive seasonal content.
//
// Example:
//
//	result, err := client.Seasons.ActiveContent(ctx)
func (s *Service) ActiveContent(ctx context.Context) (ActiveContentResponse, error) {
	client := s.base.Client
	resp, err := client.V5().GetSeasonActiveContent(client.Limiter().Wrap(ctx))
	if err != nil {
		return ActiveContentResponse{}, err
	}
	parsed, meta, err := common.Parse(resp, v5Client.ParseGetSeasonActiveContentResponse)
	if err != nil {
		return ActiveContentResponse{ResponseMeta: meta}, err
	}
	return ActiveContentResponse{Data: parsed.JSON200.Data, ResponseMeta: meta}, nil
}

type LeaderboardType = v4Client.GetSeasonLeaderboardParamsLeaderboard

const (
	LeaderboardPlayers LeaderboardType = v4Client.GetSeasonLeaderboardParamsLeaderboardPlayers
	LeaderboardTeams   LeaderboardType = v4Client.GetSeasonLeaderboardParamsLeaderboardTeams
)

type UserRankRanksResponse struct {
	Data         []v4Client.SeasonUserRankData
	ResponseMeta common.ResponseMeta
}

// UserRankById retrieves season rank data for a specific user ID.
//
// Example:
//
//	rank, err := client.Seasons.UserRankById(ctx, 12345)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("User rank: %d\n", rank.Data.Rank)
func (s *Service) UserRankById(ctx context.Context, userId int) (UserRankRanksResponse, error) {
	resp, err := s.base.Client.V4().GetSeasonUserUserIdRank(
		s.base.Client.Limiter().Wrap(ctx),
		userId,
	)
	if err != nil {
		return UserRankRanksResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSeasonUserUserIdRankResponse)
	if err != nil {
		return UserRankRanksResponse{ResponseMeta: meta}, err
	}

	return UserRankRanksResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

// UserRanksResponse contains a user’s rankings across seasons and HTTP response metadata.
type UserRanksResponse struct {
	Data         []v5Client.SeasonUserRanksEntry
	ResponseMeta common.ResponseMeta
}

// UserRanks returns the given user’s rankings across seasons.
//
// Example:
//
//	result, err := client.Seasons.UserRanks(ctx, 12345)
func (s *Service) UserRanks(ctx context.Context, userID int) (UserRanksResponse, error) {
	client := s.base.Client
	resp, err := client.V5().GetSeasonUserRanks(client.Limiter().Wrap(ctx), userID)
	if err != nil {
		return UserRanksResponse{}, err
	}
	parsed, meta, err := common.Parse(resp, v5Client.ParseGetSeasonUserRanksResponse)
	if err != nil {
		return UserRanksResponse{ResponseMeta: meta}, err
	}
	return UserRanksResponse{Data: parsed.JSON200.Data, ResponseMeta: meta}, nil
}

type LeaderboardData = v4Client.SeasonPlayersLeaderboardResponse

type LeaderboardResponse struct {
	Data         LeaderboardData
	ResponseMeta common.ResponseMeta
}

// Leaderboard retrieves season leaderboard entries for players or teams.
//
// Example:
//
//	leaderboard, err := client.Seasons.Leaderboard(ctx, seasons.LeaderboardPlayers, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Leaderboard entries: %d\n", len(leaderboard.Data.Data))
func (s *Service) Leaderboard(ctx context.Context, leaderboard LeaderboardType, params *v4Client.GetSeasonLeaderboardParams) (LeaderboardResponse, error) {
	if leaderboard == "" {
		leaderboard = LeaderboardPlayers
	}

	resp, err := s.base.Client.V4().GetSeasonLeaderboard(
		s.base.Client.Limiter().Wrap(ctx),
		v4Client.GetSeasonLeaderboardParamsLeaderboard(leaderboard),
		params,
	)
	if err != nil {
		return LeaderboardResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSeasonLeaderboardResponse)
	if err != nil {
		return LeaderboardResponse{ResponseMeta: meta}, err
	}

	return LeaderboardResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

// SeasonLeaderboardResponse contains a page of season rankings and HTTP response metadata.
type SeasonLeaderboardResponse struct {
	Data         v5Client.SeasonLeaderboardResponse
	ResponseMeta common.ResponseMeta
}

// Leaderboard returns a page of rankings for this season.
//
// Example:
//
//	result, err := client.Seasons.Season(21).Leaderboard(ctx, params)
func (h *Handle) Leaderboard(ctx context.Context, params *v5Client.GetSeasonLeaderboardParams) (SeasonLeaderboardResponse, error) {
	client := h.client
	resp, err := client.V5().GetSeasonLeaderboard(client.Limiter().Wrap(ctx), h.id, params)
	if err != nil {
		return SeasonLeaderboardResponse{}, err
	}
	parsed, meta, err := common.Parse(resp, v5Client.ParseGetSeasonLeaderboardResponse)
	if err != nil {
		return SeasonLeaderboardResponse{ResponseMeta: meta}, err
	}
	return SeasonLeaderboardResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

// LeaderboardTopResponse contains top season rankings and HTTP response metadata.
type LeaderboardTopResponse struct {
	Data         v5Client.SeasonLeaderboardTopResponse
	ResponseMeta common.ResponseMeta
}

// LeaderboardTop returns the top-ranked entries for this season.
//
// Example:
//
//	result, err := client.Seasons.Season(21).LeaderboardTop(ctx, params)
func (h *Handle) LeaderboardTop(ctx context.Context, params *v5Client.GetSeasonLeaderboardTopParams) (LeaderboardTopResponse, error) {
	client := h.client
	resp, err := client.V5().GetSeasonLeaderboardTop(client.Limiter().Wrap(ctx), h.id, params)
	if err != nil {
		return LeaderboardTopResponse{}, err
	}
	parsed, meta, err := common.Parse(resp, v5Client.ParseGetSeasonLeaderboardTopResponse)
	if err != nil {
		return LeaderboardTopResponse{ResponseMeta: meta}, err
	}
	return LeaderboardTopResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}
