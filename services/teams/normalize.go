package teams

import (
	v4client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/deref"
)

func fromAPIUserEntry(data v4client.UserEntry) UserEntry {
	return UserEntry{
		Id:          deref.Int(data.Id),
		TeamId:      deref.Int(data.TeamId),
		User:        fromAPIUser(data.User),
		UserId:      deref.Int(data.UserId),
		UserRequest: deref.Int(data.UserRequest),
	}
}

func fromAPIUser(data *v4client.User) User {
	if data == nil {
		return User{}
	}
	return User{
		AvatarThumb:      deref.String(data.AvatarThumb),
		Id:               deref.Int(data.Id),
		Name:             deref.String(data.Name),
		Points:           deref.Int(data.Points),
		RankName:         deref.String(data.RankName),
		Ranking:          userRankingToInt(data.Ranking),
		Rankings:         convert.SlicePointer(data.Rankings, fromAPIUserRanking),
		RespectedByCount: deref.Int(data.RespectedByCount),
		RootOwnsCount:    deref.Int(data.RootOwnsCount),
		UserOwnsCount:    deref.Int(data.UserOwnsCount),
	}
}

func userRankingToInt(raw *v4client.User_Ranking) int {
	if raw == nil {
		return 0
	}

	v, err := raw.AsUserRanking0()
	if err != nil {
		return 0
	}

	return v
}

func fromAPIUserRanking(data v4client.UserRanking) UserRanking {
	return UserRanking{
		Challenges: deref.Int(data.Challenges),
		CreatedAt:  deref.Time(data.CreatedAt),
		Endgame:    deref.Int(data.Endgame),
		Fc:         deref.Int(data.Fc),
		Fortress:   deref.Int(data.Fortress),
		Fr:         deref.Int(data.Fr),
		Fu:         deref.Int(data.Fu),
		Id:         deref.Int(data.Id),
		Ownership:  deref.String(data.Ownership),
		Points:     deref.Int(data.Points),
	}
}

func fromAPITeamMember(data v4client.TeamMember) TeamMember {
	return TeamMember{
		Avatar:          deref.String(data.Avatar),
		CountryCode:     deref.String(data.CountryCode),
		CountryName:     deref.String(data.CountryName),
		Id:              deref.Int(data.Id),
		Name:            deref.String(data.Name),
		Points:          deref.Int(data.Points),
		Public:          deref.Int(data.Public),
		Rank:            deref.String(data.Rank),
		RankText:        deref.String(data.RankText),
		Role:            deref.String(data.Role),
		RootBloodsCount: deref.Int(data.RootBloodsCount),
		RootOwns:        deref.Int(data.RootOwns),
		Team:            fromAPITeamMemberTeam(data.Team),
		UserBloodsCount: deref.Int(data.UserBloodsCount),
		UserOwns:        deref.Int(data.UserOwns),
	}
}

func fromAPITeamMemberTeam(data *v4client.TeamMemberTeam) TeamMemberTeam {
	if data == nil {
		return TeamMemberTeam{}
	}
	return TeamMemberTeam{
		Id:        deref.Int(data.Id),
		CaptainId: deref.Int(data.CaptainId),
	}
}

func fromAPITeamActivityItem(data v4client.TeamActivityItem) TeamActivityItem {
	return TeamActivityItem{
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
		User:              fromAPITeamActivityUser(data.User),
	}
}

func fromAPITeamActivityUser(data v4client.TeamActivityUser) TeamActivityUser {
	return TeamActivityUser{
		AvatarThumb: deref.String(data.AvatarThumb),
		Id:          data.Id,
		Name:        data.Name,
		Public:      data.Public,
	}
}
