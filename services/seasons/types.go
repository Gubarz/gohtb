package seasons

import (
	"time"

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

type SeasonActiveData struct {
	Active          bool
	Avatar          string
	Cocreators      []common.UserBasicInfoWithRespect
	Creator         []common.UserBasicInfoWithRespect
	DifficultyText  string
	Id              int
	InfoStatus      string
	Ip              string
	IsOwnedRoot     bool
	IsOwnedUser     bool
	IsReleased      bool
	IsRootBlood     bool
	IsUserBlood     bool
	MakerId         int
	Name            string
	Os              string
	PlayInfo        common.PlayInfo
	Points          int
	Poweroff        int
	Production      bool
	Release         time.Time
	ReleaseTime     time.Time
	Retired         bool
	RootBloodPoints int
	RootOwnPoints   int
	RootPoints      int
	StaticPoints    int
	Unknown         bool
	UserBloodPoints int
	UserOwnPoints   int
	UserPoints      int
}

type SeasonMachinesDataItem struct {
	Userpoints     int
	Active         bool
	Avatar         string
	DifficultyText string
	Id             int
	IsUserblood    bool
	IsOwnedRoot    bool
	IsOwnedUser    bool
	IsReleased     bool
	IsRootBlood    bool
	Name           string
	Os             string
	ReleaseTime    time.Time
	RootPoints     int
	Unknown        bool
}

type SeasonListDataItem struct {
	Active          bool
	BackgroundImage string
	EndDate         time.Time
	Id              int
	IsVisible       bool
	Logo            string
	Name            string
	StartDate       time.Time
	State           string
	Subtitle        string
}

type TopUserItem struct {
	Id         int
	Name       string
	Points     int
	UserAvatar string
}

type SeasonUserFollowerData struct {
	TopRankedFollowers []TopUserItem
	TopSeasonUsers     []TopUserItem
}

type SeasonUserRankData struct {
	FlagsToNextRank   FlagsToNextRank
	League            string
	Rank              int
	RankSuffix        string
	TotalRanks        int
	TotalSeasonPoints int
}

type FlagsToNextRank struct {
	Obtained int
	Total    int
}

type SeasonRewardsDataItem struct {
	RewardTypes SeasonRewardTypes
}

type SeasonRewardTypes struct {
	Description string
	Groups      []SeasonRewardGroupItem
	Id          int
	Name        string
	Order       float32
}

type SeasonRewardGroupItem struct {
	Description  string
	Id           int
	Image        string
	Name         string
	Order        int
	RewardTypeId int
	Rewards      []SeasonRewardItem
	Subtitle     string
}

type SeasonRewardItem struct {
	Id            int
	Image         string
	Name          string
	Order         float32
	RewardGroupId int
}
