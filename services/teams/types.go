package teams

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

type InvitationsResponse struct {
	Data         []UserEntry
	ResponseMeta common.ResponseMeta
}

type MembersResponse struct {
	Data         []TeamMember
	ResponseMeta common.ResponseMeta
}

type ActivityResponse struct {
	Data         []TeamActivityItem
	ResponseMeta common.ResponseMeta
}

type UserEntry struct {
	Id          int
	TeamId      int
	User        User
	UserId      int
	UserRequest int
}

type User struct {
	AvatarThumb      string
	Id               int
	Name             string
	Points           int
	RankName         string
	Ranking          int
	Rankings         []UserRanking
	RespectedByCount int
	RootOwnsCount    int
	UserOwnsCount    int
}

type UserRanking struct {
	Challenges int
	CreatedAt  time.Time
	Endgame    int
	Fc         int
	Fortress   int
	Fr         int
	Fu         int
	Id         int
	Ownership  string
	Points     int
	Pro        int
	Rank       int
	Respect    int
	Roots      int
	Sr         int
	Su         int
	Tr         int
	Tu         int
	UpdatedAt  time.Time
	UserId     int
	Users      int
}

type TeamMember struct {
	Avatar          string
	CountryCode     string
	CountryName     string
	Id              int
	Name            string
	Points          int
	Public          int
	Rank            int
	RankText        string
	Role            string
	RootBloodsCount int
	RootOwns        int
	Team            TeamMemberTeam
	UserBloodsCount int
	UserOwns        int
}

type TeamMemberTeam struct {
	CaptainId int
	Id        int
}

type TeamActivityItem struct {
	ChallengeCategory string
	Date              time.Time
	DateDiff          string
	FirstBlood        bool
	FlagTitle         string
	Id                int
	MachineAvatar     string
	Name              string
	ObjectType        string
	Points            int
	Type              string
	User              TeamActivityUser
}

type TeamActivityUser struct {
	AvatarThumb string
	Id          int
	Name        string
	Public      int
}
