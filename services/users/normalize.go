package users

import (
	"github.com/gubarz/gohtb/internal/deref"
	v4client "github.com/gubarz/gohtb/internal/httpclient/v4"
)

func fromAPIUserProfile(data *v4client.UserProfile) UserProfile {
	return UserProfile{
		Avatar:              deref.String(data.Avatar),
		CountryCode:         deref.String(data.CountryCode),
		CountryName:         deref.String(data.CountryName),
		CurrentRankProgress: deref.Float32(data.CurrentRankProgress),
		Github:              deref.String(data.Github),
		Id:                  deref.Int(data.Id),
		IsDedicatedVip:      deref.Bool(data.IsDedicatedVip),
		IsFollowed:          deref.Bool(data.IsFollowed),
		IsRespected:         deref.Bool(data.IsRespected),
		IsVip:               deref.Bool(data.IsVip),
		Linkedin:            deref.String(data.Linkedin),
		Name:                deref.String(data.Name),
		NextRank:            deref.String(data.NextRank),
		NextRankPoints:      deref.Float32(data.NextRankPoints),
		Points:              deref.Int(data.Points),
		Public:              deref.Bool(data.Public),
		Rank:                deref.String(data.Rank),
		RankId:              deref.Int(data.RankId),
		RankOwnership:       deref.Float32(data.RankOwnership),
		RankRequirement:     deref.Int(data.RankRequirement),
		Ranking:             deref.Int(data.Ranking),
		Respects:            deref.Int(data.Respects),
		SsoId:               deref.Bool(data.SsoId),
		SystemBloods:        deref.Int(data.SystemBloods),
		SystemOwns:          deref.Int(data.SystemOwns),
		Timezone:            deref.String(data.Timezone),
		Twitter:             deref.String(data.Twitter),
		University:          fromAPIUserProfileUniversityTeam(data.University),
		UniversityName:      deref.String(data.UniversityName),
		UserBloods:          deref.Int(data.UserBloods),
		UserOwns:            deref.Int(data.UserOwns),
		Team:                fromAPIUserProfileTeam(data.Team),
	}
}

func fromAPIUserProfileTeam(data *v4client.UserProfileTeam) UserProfileTeam {
	if data == nil {
		return UserProfileTeam{
			Id:      0,
			Name:    "",
			Ranking: 0,
			Avatar:  "",
		}
	}
	return UserProfileTeam{
		Id:      deref.Int(data.Id),
		Name:    deref.String(data.Name),
		Ranking: deref.Int(data.Ranking),
		Avatar:  deref.String(data.Avatar),
	}
}

func fromAPIUserProfileUniversityTeam(data *v4client.UserProfileUniversityTeam) UserProfileTeam {
	if data == nil {
		return UserProfileTeam{
			Id:      0,
			Name:    "",
			Ranking: 0,
			Avatar:  "",
		}
	}
	return UserProfileTeam{
		Id:      deref.Int(data.Id),
		Name:    deref.String(data.Name),
		Ranking: deref.Int(data.Ranking),
		Avatar:  deref.String(data.LogoThumbUrl),
	}
}

func fromAPIUserActivity(data v4client.UserActivityItem) UserActivityItem {
	return UserActivityItem{
		ChallengeCategory: deref.String(data.ChallengeCategory),
		Date:              data.Date,
		DateDiff:          data.DateDiff,
		FirstBlood:        data.FirstBlood,
		FlagTitle:         deref.String(data.FlagTitle),
		Id:                data.Id,
		MachineAvatar:     deref.String(data.MachineAvatar),
		Name:              data.Name,
		Points:            data.Points,
		Type:              data.Type,
		ObjectType:        data.ObjectType,
	}
}
