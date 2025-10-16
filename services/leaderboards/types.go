package leaderboards


import (
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

type Service struct {
	base service.Base
}

type Handle struct {
	client service.Client
}
// getting `rankings/users`
type RankingsUserData struct {
	Avatar_thumb string
	Challenge_bloods int
	Challenge_owns int
	Country	string
	Fortress	int
	Id	int
	Level string
	Name string
	Points int
	Rank int
	Ranks_diff int
	Root_bloods int
	Root_owns int
	User_bloods int
	User_owns int
}
type UserRankingsResponse struct {
	// Data UserRankings
	Data []RankingsUserData
	ResponseMeta common.ResponseMeta
}

// end of `rankings/users` types
