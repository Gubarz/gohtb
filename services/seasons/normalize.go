package seasons

import (
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/deref"
	v4client "github.com/gubarz/gohtb/internal/httpclient/v4"
)

func fromAPISeasonMachineActive(data *v4client.SeasonActiveData) SeasonActiveData {
	if data == nil {
		return SeasonActiveData{}
	}
	return SeasonActiveData{
		Avatar:          deref.String(data.Avatar),
		Cocreators:      convert.SlicePointer(data.Cocreators, fromAPIUserBasicInfoWithRespec),
		Creator:         convert.SlicePointer(data.Creator, fromAPIUserBasicInfoWithRespec),
		DifficultyText:  deref.String(data.DifficultyText),
		Id:              deref.Int(data.Id),
		InfoStatus:      deref.String(data.InfoStatus),
		Ip:              deref.String(data.Ip),
		IsOwnedRoot:     deref.Bool(data.IsOwnedRoot),
		IsOwnedUser:     deref.Bool(data.IsOwnedUser),
		IsReleased:      deref.Bool(data.IsReleased),
		IsRootBlood:     deref.Bool(data.IsRootBlood),
		IsUserBlood:     deref.Bool(data.IsUserBlood),
		MakerId:         deref.Int(data.MakerId),
		Name:            deref.String(data.Name),
		Os:              deref.String(data.Os),
		PlayInfo:        fromAPIPlayInfo(data.PlayInfo),
		Points:          deref.Int(data.Points),
		Poweroff:        deref.Int(data.Poweroff),
		Production:      deref.Bool(data.Production),
		Release:         deref.Time(data.Release),
		ReleaseTime:     deref.Time(data.ReleaseTime),
		Retired:         deref.Bool(data.Retired),
		RootBloodPoints: deref.Int(data.RootBloodPoints),
		RootOwnPoints:   deref.Int(data.RootOwnPoints),
		RootPoints:      deref.Int(data.RootPoints),
		StaticPoints:    deref.Int(data.StaticPoints),
		Unknown:         deref.Bool(data.Unknown),
		UserBloodPoints: deref.Int(data.UserBloodPoints),
		UserOwnPoints:   deref.Int(data.UserOwnPoints),
		UserPoints:      deref.Int(data.UserPoints),
	}
}

func fromAPIUserBasicInfoWithRespec(data v4client.UserBasicInfoWithRespect) common.UserBasicInfoWithRespect {
	return common.UserBasicInfoWithRespect{
		Id:          deref.Int(data.Id),
		Name:        deref.String(data.Name),
		Avatar:      deref.String(data.Avatar),
		IsRespected: deref.Bool(data.IsRespected),
	}
}

func fromAPIPlayInfo(data *v4client.PlayInfo) common.PlayInfo {
	if data == nil {
		return common.PlayInfo{}
	}
	return common.PlayInfo{
		ActivePlayerCount: deref.Int(data.ActivePlayerCount),
		ExpiresAt:         deref.Time(data.ExpiresAt),
		IsActive:          deref.Bool(data.IsActive),
		IsSpawned:         deref.Bool(data.IsSpawned),
		IsSpawning:        deref.Bool(data.IsSpawning),
	}
}

func fromAPISeasonListDataItem(data v4client.SeasonListDataItem) SeasonListDataItem {
	return SeasonListDataItem{
		Active:          deref.Bool(data.Active),
		BackgroundImage: deref.String(data.BackgroundImage),
		EndDate:         deref.Time(data.EndDate),
		Id:              deref.Int(data.Id),
		IsVisible:       deref.Bool(data.IsVisible),
		Logo:            deref.String(data.Logo),
		Name:            deref.String(data.Name),
		StartDate:       deref.Time(data.StartDate),
		State:           deref.String(data.State),
		Subtitle:        deref.String(data.Subtitle),
	}
}

func fromAPISeasonUserFollowerData(data *v4client.SeasonUserFollowerData) SeasonUserFollowerData {
	if data == nil {
		return SeasonUserFollowerData{}
	}
	return SeasonUserFollowerData{
		TopRankedFollowers: convert.SlicePointer(data.TopRankedFollowers, fromAPITopUserItem),
		TopSeasonUsers:     convert.SlicePointer(data.TopSeasonUsers, fromAPITopUserItem),
	}
}

func fromAPITopUserItem(data v4client.TopUserItem) TopUserItem {
	return TopUserItem{
		Id:         deref.Int(data.Id),
		Name:       deref.String(data.Name),
		Points:     deref.Int(data.Points),
		UserAvatar: deref.String(data.UserAvatar),
	}
}

func fromAPISeasonUserRankData(data *v4client.SeasonUserRankData) SeasonUserRankData {
	if data == nil {
		return SeasonUserRankData{}
	}
	flagrank, err := data.FlagsToNextRank.AsFlagsToNextRank0()
	if err != nil {
		flagrank = v4client.FlagsToNextRank0{}
	}
	return SeasonUserRankData{
		FlagsToNextRank:   fromAPIFlagsToNextRank(&flagrank),
		League:            deref.String(data.League),
		Rank:              deref.Int(data.Rank),
		RankSuffix:        deref.String(data.RankSuffix),
		TotalRanks:        deref.Int(data.TotalRanks),
		TotalSeasonPoints: deref.Int(data.TotalSeasonPoints),
	}
}

func fromAPIFlagsToNextRank(data *v4client.FlagsToNextRank0) FlagsToNextRank {
	if data == nil {
		return FlagsToNextRank{}
	}
	return FlagsToNextRank{
		Obtained: deref.Int(data.Obtained),
		Total:    deref.Int(data.Total),
	}
}

func fromAPISeasonRewardsDataItem(data v4client.SeasonRewardsDataItem) SeasonRewardsDataItem {
	return SeasonRewardsDataItem{
		RewardTypes: fromAPISeasonRewardTypes(data.RewardTypes),
	}
}

func fromAPISeasonRewardTypes(data *v4client.SeasonRewardTypes) SeasonRewardTypes {
	if data == nil {
		return SeasonRewardTypes{}
	}
	return SeasonRewardTypes{
		Description: deref.String(data.Description),
		Groups:      convert.SlicePointer(data.Groups, fromAPISeasonRewardGroupItem),
		Id:          deref.Int(data.Id),
		Name:        deref.String(data.Name),
		Order:       deref.Float32(data.Order),
	}
}

func fromAPISeasonRewardGroupItem(data v4client.SeasonRewardGroupItem) SeasonRewardGroupItem {
	return SeasonRewardGroupItem{
		Description:  deref.String(data.Description),
		Id:           deref.Int(data.Id),
		Image:        deref.String(data.Image),
		Name:         deref.String(data.Name),
		Order:        deref.Int(data.Order),
		RewardTypeId: deref.Int(data.RewardTypeId),
		Rewards:      convert.SlicePointer(data.Rewards, fromAPISeasonRewardItem),
		Subtitle:     deref.String(data.Subtitle),
	}
}

func fromAPISeasonRewardItem(data v4client.SeasonRewardItem) SeasonRewardItem {
	return SeasonRewardItem{
		Id:            deref.Int(data.Id),
		Image:         deref.String(data.Image),
		Name:          deref.String(data.Name),
		Order:         deref.Float32(data.Order),
		RewardGroupId: deref.Int(data.RewardGroupId),
	}
}

func fromAPISeasonMachinesDataItem(data v4client.SeasonMachinesDataItem) SeasonMachinesDataItem {
	return SeasonMachinesDataItem{
		Userpoints:     deref.Int(data.UserPoints),
		Active:         deref.Bool(data.Active),
		Avatar:         deref.String(data.Avatar),
		DifficultyText: deref.String(data.DifficultyText),
		Id:             deref.Int(data.Id),
		IsUserblood:    deref.Bool(data.IsUserBlood),
		IsOwnedRoot:    deref.Bool(data.IsOwnedRoot),
		IsOwnedUser:    deref.Bool(data.IsOwnedUser),
		IsReleased:     deref.Bool(data.IsReleased),
		IsRootBlood:    deref.Bool(data.IsRootBlood),
		Name:           deref.String(data.Name),
		Os:             deref.String(data.Os),
		ReleaseTime:    deref.Time(data.ReleaseTime),
		RootPoints:     deref.Int(data.RootPoints),
		Unknown:        deref.Bool(data.Unknown),
	}
}
