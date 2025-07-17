package challenges

import (
	"time"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

type Service struct {
	base    service.Base
	product string
}

type Handle struct {
	client  service.Client
	id      int
	product string
}

type ChallengeQuery struct {
	client     service.Client
	status     *v4Client.GetChallengesParamsStatus
	state      *v4Client.State
	sortBy     *v4Client.GetChallengesParamsSortBy
	sortType   *v4Client.GetChallengesParamsSortType
	difficulty *v4Client.Difficulty
	category   *v4Client.Category
	todo       *v4Client.GetChallengesParamsTodo
	page       int
	perPage    int
}

// junk remove
type CompanyItem struct {
	Id int
}

type ChallengeListResponse struct {
	Data         []ChallengeList
	ResponseMeta common.ResponseMeta
}

type ActivityResponse struct {
	Data         []ChallengeActivity
	ResponseMeta common.ResponseMeta
}

type CategoriesResponse struct {
	Data         []v4Client.Category
	ResponseMeta common.ResponseMeta
}

type WriteupResponse struct {
	Data         []CompanyItem
	ResponseMeta common.ResponseMeta
}

type WriteupOfficialResponse struct {
	Data         []CompanyItem
	ResponseMeta common.ResponseMeta
}

type ChangeLogResponse struct {
	Data         []CompanyItem
	ResponseMeta common.ResponseMeta
}

type DownloadResponse struct {
	Data         []byte
	ResponseMeta common.ResponseMeta
}

type InfoResponse struct {
	Data         Challenge
	ResponseMeta common.ResponseMeta
}

type StartResponse struct {
	Data         []CompanyItem
	ResponseMeta common.ResponseMeta
}

type StopResponse struct {
	Data         []CompanyItem
	ResponseMeta common.ResponseMeta
}

type ChallengeList struct {
	AuthUserHasReviewed bool
	Avatar              string
	CategoryId          int
	CategoryName        string
	Difficulty          string
	DifficultyChart     common.DifficultyChart
	Id                  int
	IsOwned             bool
	Name                string
	Pinned              bool
	PlayMethods         []string
	Rating              float32
	RatingCount         int
	ReleaseDate         time.Time
	Retires             ChallengeRetires
	Solves              int
	State               string
	UserDifficulty      string
}

type ChallengeRetires struct {
	name       string
	difficulty string
}

type ChallengeActivity struct {
	CreatedAt  time.Time
	Date       string
	DateDiff   string
	Type       string
	UserAvatar string
	UserId     int
	UserName   string
}

type Challenge struct {
	AuthUserHasReviewed  bool
	AuthUserSolve        bool
	AuthUserSolveTime    string
	CanAccessWalkthough  bool
	CategoryName         string
	Creator2Avatar       string
	Creator2Id           int
	Creator2Name         string
	CreatorAvatar        string
	CreatorId            int
	CreatorName          string
	Description          string
	Difficulty           string
	DifficultyChart      common.DifficultyChart
	DislikeByAuthUser    bool
	Dislikes             int
	Docker               bool
	DockerIp             string
	DockerPorts          string
	DockerStatus         string
	Download             bool
	FirstBloodTime       string
	FirstBloodUser       string
	FirstBloodUserAvatar string
	FirstBloodUserId     int
	FileName             string
	FileSize             string
	HasChangelog         bool
	Id                   int
	IsRespected          bool
	IsRespected2         bool
	IsTodo               bool
	LikeByAuthUser       bool
	Likes                int
	Name                 string
	PlayInfo             common.PlayInfoAlt
	PlayMethods          common.StringArray
	Points               int
	Recommended          int
	ReleaseDate          time.Time
	Released             int
	Retired              bool
	ReviewsCount         int
	Sha256               string
	ShowGoVip            bool
	Solves               int
	Stars                float32
	State                string
	Tags                 common.StringArray
	UserCanReview        bool
}
