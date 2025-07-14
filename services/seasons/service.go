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

// Season returns a handle for a specific season with the given ID.
// This handle can be used to perform operations related to that season,
// such as retrieving rewards, user rankings, and follower information.
func (s *Service) Season(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
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
//		fmt.Printf("Reward: %s (Points: %d)\n", reward.Name, reward.Points)
//	}
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

// UserRank retrieves the current user's ranking information for the specified season.
// This includes position, points, and other ranking details for the authenticated user.
//
// Example:
//
//	rank, err := client.Seasons.Season(123).UserRank(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Current rank: %d (Points: %d)\n", rank.Data.Position, rank.Data.Points)
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

// UserFollowers retrieves follower information for the current user in the specified season.
// This includes details about users following the authenticated user during the season.
//
// Example:
//
//	followers, err := client.Seasons.Season(123).UserFollowers(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Followers: %d\n", len(followers.Data.Followers))
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
//		fmt.Printf("Machine: %s (Difficulty: %s)\n", machine.Name, machine.Difficulty)
//	}
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
