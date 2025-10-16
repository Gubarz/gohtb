package users

import (
	"time"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
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

type UserProfile struct{ v4Client.UserProfile }
type UserProfileTeam = v4Client.UserProfileTeam
type UserActivityItem = v4Client.UserActivityItem
type UserActivity = []UserActivityItem

func wrapUserProfile(x v4Client.UserProfile) UserProfile { return UserProfile{x} }
