package machines

import (
	"time"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	v5Client "github.com/gubarz/gohtb/httpclient/v5"
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

type InfoResponse struct {
	Data         MachineProfileInfo
	ResponseMeta common.ResponseMeta
}

type UnreleasedDataItems []MachineUnreleasedData

type MachineUnreleasedResponse struct {
	Data         UnreleasedDataItems
	ResponseMeta common.ResponseMeta
}

type MachineDataItems []MachineData

type MachinePaginatedResponse struct {
	Data         MachineDataItems
	Pagination   PagingMeta
	ResponseMeta common.ResponseMeta
}

type MachineResponse struct {
	Data         MachineDataItems
	ResponseMeta common.ResponseMeta
}

type OwnResponse struct {
	Data         MachineOwnResponse
	ResponseMeta common.ResponseMeta
}

type PagingMeta struct {
	CurrentPage int
	PerPage     int
	Total       int
	TotalPages  int
	Count       int
}

type ActiveResponse struct {
	Data         ActiveMachineInfo
	ResponseMeta common.ResponseMeta
}

type ActiveQuery struct {
	client        service.Client
	perPage       int
	page          int
	showCompleted string
	sortBy        v4Client.GetMachinePaginatedParamsSortBy
	sortType      v4Client.GetMachinePaginatedParamsSortType
	difficulty    v4Client.Difficulty
	os            v4Client.Os
	keyword       v4Client.Keyword
}

type RetiredQuery struct {
	client        service.Client
	perPage       int
	page          int
	showCompleted string
	sortBy        v4Client.GetMachinePaginatedParamsSortBy
	sortType      v4Client.GetMachinePaginatedParamsSortType
	difficulty    v4Client.Difficulty
	os            v4Client.Os
	keyword       v4Client.Keyword
	tags          v4Client.Tags
}

type UnreleasedQuery struct {
	client     service.Client
	perPage    int
	page       int
	difficulty v4Client.Difficulty
	os         v4Client.Os
	keyword    v4Client.Keyword
}

type Machine struct {
	AvatarThumbUrl string
	Id             int
	Name           string
	Os             string
}

type MachineAcitivtyItem struct {
	BloodType  string
	CreatedAt  string
	Date       string
	DateDiff   string
	Type       string
	UserAvatar string
	UserId     int
	UserName   string
}

type MachineActivityIdResponse struct {
	Info MachineActivityInfo
}

type MachineActivityInfo struct {
	Activity []MachineAcitivtyItem
	Server   string
}

type MachineAttackDetails struct {
	MachineAttackPaths common.TeamMachineAttackPaths
	MachineOwns        MachineOwnsCard
}

type MachineAttackPathItem struct {
	AvgUniversitysSolved float32
	Name                 string
	Solved               float32
	Total                float32
}

type MachineAttackPaths struct {
	BinaryAnalysis            MachineAttackPathItem
	BinaryExploitation        MachineAttackPathItem
	ConfigurationAnalysis     MachineAttackPathItem
	Fuzzing                   MachineAttackPathItem
	Impersonation             MachineAttackPathItem
	PacketCaptureAnalysis     MachineAttackPathItem
	Pivoting                  MachineAttackPathItem
	Reconnaissance            MachineAttackPathItem
	UserEnumeration           MachineAttackPathItem
	WebSiteStructureDiscovery MachineAttackPathItem
}

type MachineCard1 struct {
	Avatar         string
	DifficultyText string
	Id             int
	IsTodo         bool
	Maker          common.UserBasicInfo
	Maker2         common.UserBasicInfo
	Name           string
	Os             string
	Points         int
	Release        time.Time
	Retired        int
	RetiredId      int
	TypeCard       string
}

type MachineCard2 struct {
	AuthUserInRootOwns bool
	AuthUserInUserOwns bool
	Avatar             string
	DifficultyText     string
	FeedbackForChart   common.DifficultyChart
	Id                 int
	Name               string
	Os                 string
	Points             int
	Release            time.Time
	Retired            int
	RetiredDate        time.Time
	RootOwnsCount      int
	Stars              float32
	TypeCard           string
	UserOwnsCount      int
}

type MachineChangeLogItem struct {
	CreatedAt   string
	Description string
	Id          int
	MachineId   int
	Released    int
	Title       string
	Type        string
	UpdatedAt   string
	UserId      int
}

type MachineChangelogIdResponse struct {
	Info []MachineChangeLogItem
}

type MachineGraphActivityIdResponse struct {
	Info MachineGraphActivityInfo
}

type MachineGraphActivityInfo struct {
	Periods    [][]string
	Resets     []int
	SystemOwns []int
	UserOwns   []int
}

type MachineGraphMatrixIdResponse struct {
	Info MachineGraphMatrixInfo
}

type MachineGraphMatrixInfo struct {
	Aggregate common.MatrixInfo
	Maker     common.MatrixInfo
	User      common.MatrixInfo
}

type MachineGraphOwnsDifficultyIdResponse struct {
	Info MachineGraphOwnsDifficultyInfoItems
}

type MachineGraphOwnsDifficultyInfoItem struct {
	Root int
	User int
}

type MachineGraphOwnsDifficultyInfoItems struct {
	N1  MachineGraphOwnsDifficultyInfoItem
	N10 MachineGraphOwnsDifficultyInfoItem
	N2  MachineGraphOwnsDifficultyInfoItem
	N3  MachineGraphOwnsDifficultyInfoItem
	N4  MachineGraphOwnsDifficultyInfoItem
	N5  MachineGraphOwnsDifficultyInfoItem
	N6  MachineGraphOwnsDifficultyInfoItem
	N7  MachineGraphOwnsDifficultyInfoItem
	N8  MachineGraphOwnsDifficultyInfoItem
	N9  MachineGraphOwnsDifficultyInfoItem
}

type MachineOwns struct {
	Solved float32
	Total  float32
}

type MachineOwnsCard struct {
	Solved int
	Total  int
}

type MachineOwnsTabloid struct {
	AvatarThumpUrl string
	AvatarUrl      string
	Id             int
	Name           string
}

type MachineOwnsTopIdResponse struct {
	Info []MachineOwnsTopItem
}

type MachineOwnsTopItem struct {
	Avatar      string
	Id          int
	IsRootBlood bool
	IsUserBlood bool
	Name        string
	OwnDate     string
	Position    int
	RankId      int
	RankText    string
	RootOwnTime string
	UserOwnDate string
	UserOwnTime string
}

type PlayInfoAlt struct {
	ExpiresAt time.Time
	Ip        string
	Ports     []int
	Status    string
}

type MachineProfileResponse struct {
	Info MachineProfileInfo
}

type MachineRecommendedRetiredCard struct {
	Avatar           string
	DifficultyText   string
	FeedbackForChart common.DifficultyChart
	Id               int
	Name             string
	Os               string
	Release          string
	Retired          int
	RetiredDate      string
}

type MachineRecommendedRetiredResponse struct {
	Card1 MachineRecommendedRetiredCard
	Card2 MachineRecommendedRetiredCard
}

type MachineReview struct {
	Id       int
	Reviewed bool
}

type MachineReviewRequest struct {
	Headline string
	Id       int
	Review   string
	Stars    float32
}

type MachineReviewResponse struct {
	Message []MachineReview
}

type MachineReviewsMessageItem struct {
	AuthUserInHelpfulReviews bool
	CreatedAt                string
	Difficulty               int
	Featured                 int
	HelpfulReviews           []common.HelpfulReviews
	HelpfulReviewsCount      float32
	Id                       int
	Message                  string
	Released                 float32
	Stars                    float32
	Title                    string
	User                     common.UserBasicInfo
	UserId                   int
}

type MachineReviewsResponse struct {
	Average float32
	Count   float32
	Message []MachineReviewsMessageItem
}

type MachineReviewsUserIdResponse struct {
	Message string
}

type MachineTagIdResponse struct {
	Info []MachineTagItems
}

type MachineTagItems struct {
	Category string
	Id       int
	Name     string
}

type MachineTagsListResponse struct {
	Info []common.TagCategory
}

type MachineTasksData_Item any

type MachineTasksResponse struct {
	Data []MachineTasksData_Item
}

type MachineWalkthroughIdResponse struct {
	Message MachineWalkthroughMessage
}

type MachineWalkthroughMessage struct {
	Official    MachineWalkthroughMessageOfficial
	UnderReview string
	Video       MachineWalkthroughVideo
	Writeups    []MachineWalkthroughMessageWriteupsItem
}

type MachineWalkthroughMessageOfficial struct {
	DislikedByUser             float32
	Filename                   string
	LikedByUser                bool
	Rating                     float32
	Sha256                     string
	ThresholdForDisplayReached int
	TotalRatings               int
}

type MachineWalkthroughMessageWriteupsItem struct {
	CreatedAt      string
	DislikedByUser string
	Id             int
	LanguageCode   string
	LanguageName   string
	LikedByUser    string
	Rating         int
	TotalRatings   int
	Url            string
	UserAvatar     string
	UserId         int
	UserName       string
}

type MachineWalkthroughOfficialFeedbackChoicesResponse struct {
	FeedbackChoices []string
}

type MachineWalkthroughRandomResponse struct {
	Message common.UserBasicInfo
}

type MachineWalkthroughVideo struct {
	CreatorAvatar string
	CreatorId     int
	CreatorName   string
	YoutubeId     string
}

type MachineWalkthroughsLanguageListItem struct {
	FullName  string
	Id        int
	ShortName string
}

type MachineWalkthroughsLanguageListResponse struct {
	Info []MachineWalkthroughsLanguageListItem
}

type Credentials struct {
	Username string
	Password string
}

type MachineData = v4Client.MachineData
type MachinePlayInfo = v4Client.MachinePlayInfo
type MachineRetiring = v4Client.MachineRetiring
type MachineUnreleasedData = v4Client.MachineUnreleasedData
type PlayInfo = v4Client.PlayInfo
type Label = v4Client.Label
type ActiveMachineInfo = v4Client.ActiveMachineInfo
type MachineOwnResponse = v5Client.MachineOwnResponse
type RankUpdate = v5Client.RankUpdate
type NewRank = v5Client.NewRank

type MachineProfileInfo struct {
	v4Client.MachineProfileInfo
	Credentials
	IsAssumedBreach bool
}

func wrapMachineProfileInfo(x v4Client.MachineProfileInfo) MachineProfileInfo {
	return MachineProfileInfo{MachineProfileInfo: x}
}
