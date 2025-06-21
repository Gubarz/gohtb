package machines

import (
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/deref"
	v4Client "github.com/gubarz/gohtb/internal/httpclient/v4"
	v5Client "github.com/gubarz/gohtb/internal/httpclient/v5"
)

func fromAPIMachineData(data v4Client.MachineData) MachineData {
	return MachineData{
		Active:              deref.Bool(data.Active),
		AuthUserHasReviewed: deref.Bool(data.AuthUserHasReviewed),
		AuthUserInRootOwns:  deref.Bool(data.AuthUserInRootOwns),
		AuthUserInUserOwns:  deref.Bool(data.AuthUserInUserOwns),
		Avatar:              deref.String(data.Avatar),
		Difficulty:          deref.Int(data.Difficulty),
		DifficultyText:      deref.String(data.DifficultyText),
		EasyMonth:           deref.Int(data.EasyMonth),
		FeedbackForChart:    common.FromAPIDifficultyChart(data.FeedbackForChart),
		Free:                deref.Bool(data.Free),
		Id:                  deref.Int(data.Id),
		Ip:                  deref.String(data.Ip),
		IsTodo:              isTodo(data.IsTodo),
		IsCompetitive:       deref.Bool(data.IsCompetitive),
		Labels:              convert.Slice(*data.Labels, fromAPILabel),
		Name:                deref.String(data.Name),
		Os:                  deref.String(data.Os),
		PlayInfo:            fromAPIMachinePlayInfo(*data.PlayInfo),
		Points:              deref.Int(data.Points),
		Poweroff:            deref.Int(data.Poweroff),
		Recommended:         deref.Int(data.Recommended),
		Release:             deref.Time(data.Release),
		RootOwnsCount:       deref.Int(data.RootOwnsCount),
		SpFlag:              deref.Int(data.SpFlag),
		Star:                deref.Float32(data.Star),
		StaticPoints:        deref.Int(data.StaticPoints),
		UserOwnsCount:       deref.Int(data.UserOwnsCount),
	}
}

func isTodo(v *v4Client.MachineData_IsTodo) bool {
	if v == nil {
		return false
	}

	if b, err := v.AsMachineDataIsTodo0(); err == nil {
		return b
	}
	if _, err := v.AsMachineDataIsTodo1(); err == nil {
		return true
	}
	return false
}

func fromAPIMachinePlayInfo(data v4Client.MachinePlayInfo) MachinePlayInfo {
	return MachinePlayInfo{
		ExpiresAt: deref.String(data.ExpiresAt),
		IsActive:  deref.Bool(data.IsActive),
	}
}

func fromAPIMachineProfileInfo(data v4Client.MachineProfileInfo) MachineProfileInfo {
	return MachineProfileInfo{
		AcademyModules:             convert.Slice(*data.AcademyModules, common.FromAPIAcademyModule),
		Active:                     deref.Bool(data.Active),
		AuthUserFirstRootTime:      deref.String(data.AuthUserFirstRootTime),
		AuthUserFirstUserTime:      deref.String(data.AuthUserFirstUserTime),
		AuthUserHasReviewed:        deref.Bool(data.AuthUserHasReviewed),
		AuthUserHasSubmittedMatrix: deref.Bool(data.AuthUserHasSubmittedMatrix),
		AuthUserInRootOwns:         deref.Bool(data.AuthUserInRootOwns),
		AuthUserInUserOwns:         deref.Bool(data.AuthUserInUserOwns),
		Avatar:                     deref.String(data.Avatar),
		CanAccessWalkthrough:       deref.Bool(data.CanAccessWalkthrough),
		DifficultyText:             deref.String(data.DifficultyText),
		FeedbackForChart:           common.FromAPIDifficultyChart(data.FeedbackForChart),
		Free:                       deref.Bool(data.Free),
		HasChangelog:               deref.Bool(data.HasChangelog),
		Id:                         deref.Int(data.Id),
		InfoStatus:                 deref.String(data.InfoStatus),
		Ip:                         deref.String(data.Ip),
		IsGuidedEnabled:            deref.Bool(data.IsGuidedEnabled),
		IsTodo:                     deref.Bool(data.IsTodo),
		MachineMode:                deref.String(data.MachineMode),
		Maker:                      common.FromAPIMaker(data.Maker),
		Maker2:                     common.FromAPIMaker(data.Maker2),
		Name:                       deref.String(data.Name),
		Os:                         deref.String(data.Os),
		OwnRank:                    deref.Int(data.OwnRank),
		PlayInfo:                   fromAPIPlayInfoAlt(data.PlayInfo),
		Points:                     deref.Int(data.Points),
		Recommended:                deref.Bool(data.Recommended),
		Release:                    deref.Time(data.Release),
		Retired:                    deref.Bool(data.Retired),
		ReviewsCount:               deref.Int(data.ReviewsCount),
		RootBlood:                  common.FromAPIBloodInfo(*data.RootBlood),
		RootOwnsCount:              deref.Int(data.RootOwnsCount),
		SeasonId:                   deref.Int(data.SeasonId),
		ShowGoVip:                  deref.Bool(data.ShowGoVip),
		ShowGoVipServer:            deref.Bool(data.ShowGoVipServer),
		SpFlag:                     deref.Int(data.SpFlag),
		Stars:                      deref.Float32(data.Stars),
		StartMode:                  deref.String(data.StartMode),
		StaticPoints:               deref.Int(data.StaticPoints),
		Synopsis:                   deref.String(data.Synopsis),
		UserBlood:                  common.FromAPIBloodInfo(*data.UserBlood),
		UserCanReview:              deref.Bool(data.UserCanReview),
		UserOwnsCount:              deref.Int(data.UserOwnsCount),
	}
}

func fromAPIMachineRetiring(data v4Client.MachineRetiring) MachineRetiring {
	return MachineRetiring{
		Avatar:         deref.String(data.Avatar),
		DifficultyText: deref.String(data.DifficultyText),
		Id:             deref.Int(data.Id),
		Name:           deref.String(data.Name),
		Os:             deref.String(data.Os),
	}
}

func fromAPIMachineUnreleasedData(data v4Client.MachineUnreleasedData) MachineUnreleasedData {
	return MachineUnreleasedData{
		Avatar:         deref.String(data.Avatar),
		CoCreators:     convert.Slice(*data.CoCreators, common.FromAPIUserBasicInfo),
		Difficulty:     deref.Int(data.Difficulty),
		DifficultyText: deref.String(data.DifficultyText),
		FirstCreator:   convert.Slice(*data.FirstCreator, common.FromAPIUserBasicInfo),
		Id:             deref.Int(data.Id),
		Name:           deref.String(data.Name),
		Os:             deref.String(data.Os),
		Release:        deref.String(data.Release),
		Retiring:       fromAPIMachineRetiring(*data.Retiring),
	}
}

func fromAPIPlayInfoAlt(data *v4Client.PlayInfoAlt) PlayInfoAlt {
	return PlayInfoAlt{
		ExpiresAt: deref.Time(data.ExpiresAt),
		Ip:        deref.String(data.Ip),
		Ports:     deref.Slice(data.Ports),
		Status:    deref.String(data.Status),
	}
}

func fromAPILabel(data v4Client.Label) Label {
	return Label{
		Color: deref.String(data.Color),
		Name:  deref.String(data.Name),
	}
}

func fromAPIActiveMachineInfo(data v4Client.ActiveMachineInfoWrapper) ActiveMachineInfo {
	return ActiveMachineInfo{
		Avatar:      deref.String(data.Avatar),
		ExpiresAt:   deref.String(data.ExpiresAt),
		Id:          deref.Int(data.Id),
		Ip:          deref.String(data.Ip),
		IsSpawning:  deref.Bool(data.IsSpawning),
		LabServer:   deref.String(data.LabServer),
		Name:        deref.String(data.Name),
		TierId:      deref.String(data.TierId),
		Type:        deref.String(data.Type),
		Voted:       deref.String(data.Voted),
		Voting:      deref.String(data.Voting),
		VpnServerId: deref.Int(data.VpnServerId),
	}
}

func fromAPIMachineOwnResponse(data v5Client.MachineOwnResponse) MachineOwnResponse {
	return MachineOwnResponse{
		BloodPoints:      data.BloodPoints,
		BloodTaken:       data.BloodTaken,
		Id:               data.Id,
		IsStartingPoint:  data.IsStartingPoint,
		LeagueRank:       fromAPIRankUpdate(data.LeagueRank),
		MachineCompleted: data.MachineCompleted,
		MachinePwned:     data.MachinePwned,
		MachineState:     data.MachineState,
		Message:          data.Message,
		OwnType:          string(data.OwnType),
		Points:           data.Points,
		Status:           data.Status,
		Success:          data.Success,
		UserRank:         fromAPIRankUpdate(data.UserRank),
	}
}

func fromAPIRankUpdate(data *v5Client.RankUpdate) RankUpdate {
	return RankUpdate{
		Changed: deref.Bool(data.Changed),
		NewRank: fromAPINewRank(data.NewRank),
	}
}

func fromAPINewRank(data *v5Client.NewRank) NewRank {
	return NewRank{
		Id:   deref.Int(data.Id),
		Text: deref.String(data.Text),
	}
}
