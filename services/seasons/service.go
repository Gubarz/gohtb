package seasons

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

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

type SeasonRewardsDataItem = v4Client.SeasonRewardsDataItem

type RewardsResponse struct {
	Data         []SeasonRewardsDataItem
	ResponseMeta common.ResponseMeta
}

// Rewards retrieves the rewards available for the specified season.
// This includes information about prizes, achievements, and other rewards
// that can be earned during the season.
//
// Example:
//
//	rewards, err := client.Seasons.Season(123).Rewards(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, reward := range rewards.Data {
//		fmt.Printf("Reward: %s\n", reward.RewardTypes.Name)
//	}
func (h *Handle) Rewards(ctx context.Context) (RewardsResponse, error) {
	resp, err := h.client.V4().GetSeasonRewards(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return RewardsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSeasonRewardsResponse)
	if err != nil {
		return RewardsResponse{ResponseMeta: meta}, err
	}

	return RewardsResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type SeasonUserRankData = v4Client.SeasonUserRankData

type UserRankResponse struct {
	Data         SeasonUserRankData
	ResponseMeta common.ResponseMeta
}

// UserRank retrieves the current user's ranking information for the specified season.
// This includes position, points, and other ranking details for the authenticated user.
//
// Example:
//
//	rank, err := client.Seasons.Season(123).UserRank(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Current rank: %d (Points: %d)\n", rank.Data.Rank, rank.Data.TotalSeasonPoints)
func (h *Handle) UserRank(ctx context.Context) (UserRankResponse, error) {
	resp, err := h.client.V4().GetSeasonUserRank(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return UserRankResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSeasonUserRankResponse)
	if err != nil {
		return UserRankResponse{ResponseMeta: meta}, err
	}

	return UserRankResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type SeasonUserFollowerData = v4Client.SeasonUserFollowerData

type UserFollowersResponse struct {
	Data         SeasonUserFollowerData
	ResponseMeta common.ResponseMeta
}

// UserFollowers retrieves follower information for the current user in the specified season.
// This includes details about users following the authenticated user during the season.
//
// Example:
//
//	followers, err := client.Seasons.Season(123).UserFollowers(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Top ranked followers: %d\n", len(followers.Data.TopRankedFollowers))
func (h *Handle) UserFollowers(ctx context.Context) (UserFollowersResponse, error) {
	resp, err := h.client.V4().GetSeasonUserFollowers(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return UserFollowersResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSeasonUserFollowersResponse)
	if err != nil {
		return UserFollowersResponse{ResponseMeta: meta}, err
	}

	return UserFollowersResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
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
//	machines, err := client.Seasons.Machines(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, machine := range machines.Data {
//		fmt.Printf("Machine: %s (Difficulty: %s)\n", machine.Name, machine.DifficultyText)
//	}
func (s *Service) Machines(ctx context.Context) (MachinesResponse, error) {
	resp, err := s.base.Client.V4().GetSeasonMachines(s.base.Client.Limiter().Wrap(ctx))
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

type LeaderboardType = v4Client.GetSeasonLeaderboardParamsLeaderboard

const (
	LeaderboardPlayers LeaderboardType = v4Client.GetSeasonLeaderboardParamsLeaderboardPlayers
	LeaderboardTeams   LeaderboardType = v4Client.GetSeasonLeaderboardParamsLeaderboardTeams
)

type LeaderboardTopType = v4Client.GetSeasonLeaderboardTopParamsLeaderboard

const (
	LeaderboardTopPlayers LeaderboardTopType = v4Client.GetSeasonLeaderboardTopParamsLeaderboardPlayers
	LeaderboardTopTeams   LeaderboardTopType = v4Client.GetSeasonLeaderboardTopParamsLeaderboardTeams
)

type LeaderboardTopPeriod = v4Client.GetSeasonLeaderboardTopParamsPeriod

const (
	LeaderboardTopPeriod1W LeaderboardTopPeriod = v4Client.GetSeasonLeaderboardTopParamsPeriodN1W
	LeaderboardTopPeriod1M LeaderboardTopPeriod = v4Client.GetSeasonLeaderboardTopParamsPeriodN1M
	LeaderboardTopPeriod3M LeaderboardTopPeriod = v4Client.GetSeasonLeaderboardTopParamsPeriodN3M
	LeaderboardTopPeriod6M LeaderboardTopPeriod = v4Client.GetSeasonLeaderboardTopParamsPeriodN6M
	LeaderboardTopPeriod1Y LeaderboardTopPeriod = v4Client.GetSeasonLeaderboardTopParamsPeriodN1Y
)

type EndData = v4Client.SeasonEndData

type EndResponse struct {
	Data         EndData
	ResponseMeta common.ResponseMeta
}

// End retrieves season-end summary data for a specific user in this season.
//
// Example:
//
//	end, err := client.Seasons.Season(3).End(ctx, 12345)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Season-end rank: %d\n", end.Data.Rank.Current)
func (h *Handle) End(ctx context.Context, userId int) (EndResponse, error) {
	resp, err := h.client.V4().GetSeasonEnd(
		h.client.Limiter().Wrap(ctx),
		h.id,
		userId,
	)
	if err != nil {
		return EndResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSeasonEndResponse)
	if err != nil {
		return EndResponse{ResponseMeta: meta}, err
	}

	return EndResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type CompletedMachinesData = v4Client.SeasonCompletedMachineData

type CompletedMachinesResponse struct {
	Data         CompletedMachinesData
	ResponseMeta common.ResponseMeta
}

// MachinesCompleted retrieves completed machine counters for this season.
//
// Example:
//
//	completed, err := client.Seasons.Season(3).MachinesCompleted(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Season flags: %d\n", completed.Data.SeasonFlags)
func (h *Handle) MachinesCompleted(ctx context.Context) (CompletedMachinesResponse, error) {
	resp, err := h.client.V4().GetSeasonMachinesCompleted(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return CompletedMachinesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSeasonMachinesCompletedResponse)
	if err != nil {
		return CompletedMachinesResponse{ResponseMeta: meta}, err
	}

	return CompletedMachinesResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type UserRankRanksResponse struct {
	Data         []SeasonUserRankData
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

type LeaderboardTopData = v4Client.SeasonPlayersLeaderboardTopResponse

type LeaderboardTopResponse struct {
	Data         LeaderboardTopData
	ResponseMeta common.ResponseMeta
}

// LeaderboardTop retrieves top season leaderboard entries for this season.
//
// Example:
//
//	top, err := client.Seasons.Season(3).LeaderboardTop(ctx, seasons.LeaderboardTopPlayers, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Top entries: %d\n", len(top.Data.Data))
func (h *Handle) LeaderboardTop(ctx context.Context, leaderboard LeaderboardTopType, params *v4Client.GetSeasonLeaderboardTopParams) (LeaderboardTopResponse, error) {
	if leaderboard == "" {
		leaderboard = LeaderboardTopPlayers
	}
	if params == nil {
		params = &v4Client.GetSeasonLeaderboardTopParams{
			Period: v4Client.GetSeasonLeaderboardTopParamsPeriod(LeaderboardTopPeriod1M),
		}
	}

	resp, err := h.client.V4().GetSeasonLeaderboardTop(
		h.client.Limiter().Wrap(ctx),
		v4Client.GetSeasonLeaderboardTopParamsLeaderboard(leaderboard),
		h.id,
		params,
	)
	if err != nil {
		return LeaderboardTopResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSeasonLeaderboardTopResponse)
	if err != nil {
		return LeaderboardTopResponse{ResponseMeta: meta}, err
	}

	return LeaderboardTopResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}
