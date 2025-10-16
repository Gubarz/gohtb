package leaderboards


import (

	v4client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/deref"
)

func fromLeaderBoardDataUsers(data v4client.RankingsUserData) RankingsUserData {
	// for `rankings/users`
	return RankingsUserData{
		Avatar_thumb: deref.String(data.AvatarThumb),
		Challenge_bloods: deref.Int(data.ChallengeBloods),
		Challenge_owns: deref.Int(data.ChallengeOwns),
		Country: deref.String(data.Country),
		Fortress: deref.Int(data.Fortress),
		Id: deref.Int(data.Id),
		Level: deref.String(data.Level),
		Name: deref.String(data.Name),
		Points: deref.Int(data.Points),
		Rank: deref.Int(data.Rank),
		Ranks_diff: deref.Int(data.RanksDiff),
		Root_bloods: deref.Int(data.RootBloods),
		Root_owns: deref.Int(data.RootOwns),
		User_bloods: deref.Int(data.UserBloods),
		User_owns: deref.Int(data.UserOwns),
	}
}



func fromLeaderBoardDataTeams(data v4client.RankingsTeamItem) RankingsTeamItem {
	// for `rankings/users`
	return RankingsTeamItem{
		Avatar_thumb: deref.String(data.AvatarThumbUrl),
		Challenge_bloods: deref.Int(data.ChallengeBloods),
		Challenge_owns: deref.Int(data.ChallengeOwns),
		Country: deref.String(data.Country),
		Fortress: deref.Int(data.Fortress),
		Id: deref.Int(data.Id),
		Name: deref.String(data.Name),
		Points: deref.Int(data.Points),
		Rank: deref.Int(data.Rank),
		Ranks_diff: deref.Int(data.RanksDiff),
		Root_bloods: deref.Int(data.RootBloods),
		Root_owns: deref.Int(data.RootOwns),
		User_bloods: deref.Int(data.UserBloods),
		User_owns: deref.Int(data.UserOwns),
	}
}


func fromLeaderBoardDataCountries(data v4client.RankingsCountriesItem) RankingsCountries{
	return RankingsCountries{
		Challenge_owns: deref.Int(data.ChallengeOwns),
		Country: deref.String(data.Country),
		Fortress: deref.Int(data.Fortress),
		Members: deref.Int(data.Members),
		Name: deref.String(data.Name),
		Points: deref.Int(data.Points),
		Rank: deref.Int(data.Rank),
		Ranks_diff: deref.Int(data.RanksDiff),
		Root_bloods: deref.Int(data.RootBloods),
		Root_owns: deref.Int(data.RootOwns),
		User_bloods: deref.Int(data.UserBloods),
		User_owns: deref.Int(data.UserOwns),
	}

}