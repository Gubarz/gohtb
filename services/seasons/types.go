package seasons

import (
	v4client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

type Service struct {
	base service.Base
}

type Handle struct {
	client service.Client
	id     int
}

type ActiveMachineResponse struct {
	Data         SeasonActiveData
	ResponseMeta common.ResponseMeta
}

type MachinesResponse struct {
	Data         []SeasonMachinesDataItem
	ResponseMeta common.ResponseMeta
}

type ListResponse struct {
	Data         []SeasonListDataItem
	ResponseMeta common.ResponseMeta
}

type UserFollowersResponse struct {
	Data         SeasonUserFollowerData
	ResponseMeta common.ResponseMeta
}

type UserRankResponse struct {
	Data         SeasonUserRankData
	ResponseMeta common.ResponseMeta
}

type RewardsResponse struct {
	Data         []SeasonRewardsDataItem
	ResponseMeta common.ResponseMeta
}

type SeasonActiveData = v4client.SeasonActiveData
type SeasonListDataItem = v4client.SeasonListDataItem
type SeasonUserFollowerData = v4client.SeasonUserFollowerData
type SeasonUserRankData = v4client.SeasonUserRankData
type FlagsToNextRank = v4client.FlagsToNextRank0
type SeasonRewardsDataItem = v4client.SeasonRewardsDataItem
type SeasonRewardTypes = v4client.SeasonRewardTypes
type SeasonRewardGroupItem = v4client.SeasonRewardGroupItem
type SeasonRewardItem = v4client.SeasonRewardItem
type SeasonMachinesDataItem = v4client.SeasonMachinesDataItem
type TopUserItem = v4client.TopUserItem
