package users

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

type ProfileActivityResposnse struct {
	Data         []UserActivityItem
	ResponseMeta common.ResponseMeta
}

type ProfileBasicResponse struct {
	Data         UserProfile
	ResponseMeta common.ResponseMeta
}

type UserActivityItemChallenge struct {
	ChallengeCategory string
	Date              time.Time
	DateDiff          string
	FirstBlood        bool
	Id                int
	Name              string
	Points            int
	Type              string
}

type UserActivityItemFortress struct {
	Date       time.Time
	DateDiff   string
	FirstBlood bool
	FlagTitle  string
	Id         int
	Name       string
	Points     int
	Type       string
}

type UserActivityItemMachine struct {
	Date          time.Time
	DateDiff      string
	FirstBlood    bool
	Id            int
	MachineAvatar string
	Name          string
	Points        int
	Type          string
}

type UserActivityItem struct {
	ChallengeCategory string
	Date              time.Time
	DateDiff          string
	FirstBlood        bool
	FlagTitle         string
	Id                int
	MachineAvatar     string
	Name              string
	Points            int
	Type              string
	ActivityType      string
	ObjectType        string
}

type UserProfile struct {
	Avatar              string
	CountryCode         string
	CountryName         string
	CurrentRankProgress float32
	Description         string
	Github              string
	Id                  int
	IsDedicatedVip      bool
	IsFollowed          bool
	IsRespected         bool
	IsVip               bool
	Linkedin            string
	Name                string
	NextRank            string
	NextRankPoints      float32
	Points              int
	Public              bool
	Rank                string
	RankId              int
	RankOwnership       float32
	RankRequirement     int
	Ranking             int
	Respects            int
	SsoId               bool
	SystemBloods        int
	SystemOwns          int
	Team                UserProfileTeam
	Timezone            string
	Twitter             string
	University          string
	UniversityName      string
	UserBloods          int
	UserOwns            int
}

type UserProfileTeam struct {
	Avatar  string
	Id      int
	Name    string
	Ranking int
}
