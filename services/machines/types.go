package machines

import (
	"time"

	"github.com/gubarz/gohtb/internal/common"
	v4Client "github.com/gubarz/gohtb/internal/httpclient/v4"
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

type MachineOwnResponse struct {
	BloodPoints      int
	BloodTaken       int
	Id               int
	IsStartingPoint  bool
	LeagueRank       RankUpdate
	MachineCompleted bool
	MachinePwned     bool
	MachineState     string
	Message          string
	OwnType          string
	Points           int
	Status           int
	Success          bool
	UserRank         RankUpdate
}

type NewRank struct {
	Id   int
	Text string
}

type RankUpdate struct {
	Changed bool
	NewRank NewRank
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
	sortBy        *v4Client.GetMachinePaginatedParamsSortBy
	sortType      *v4Client.GetMachinePaginatedParamsSortType
	difficulty    v4Client.Difficulty
	os            v4Client.Os
	keyword       *v4Client.Keyword
}

type RetiredQuery struct {
	client        service.Client
	perPage       int
	page          int
	showCompleted string
	sortBy        *v4Client.GetMachinePaginatedParamsSortBy
	sortType      *v4Client.GetMachinePaginatedParamsSortType
	difficulty    v4Client.Difficulty
	os            v4Client.Os
	keyword       *v4Client.Keyword
	tags          v4Client.Tags
}

type UnreleasedQuery struct {
	client        service.Client
	perPage       int
	page          int
	showCompleted string
	sortBy        *v4Client.GetMachinePaginatedParamsSortBy
	sortType      *v4Client.GetMachinePaginatedParamsSortType
	difficulty    v4Client.Difficulty
	os            v4Client.Os
	keyword       *v4Client.Keyword
}

type ActiveMachineInfo struct {
	Avatar      string
	ExpiresAt   string
	Id          int
	Ip          string
	IsSpawning  bool
	LabServer   string
	Name        string
	TierId      string
	Type        string
	Voted       string
	Voting      string
	VpnServerId int
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

type Label struct {
	Color string
	Name  string
}

type MachineData struct {
	Active              bool
	AuthUserHasReviewed bool
	AuthUserInRootOwns  bool
	AuthUserInUserOwns  bool
	Avatar              string
	Difficulty          int
	DifficultyText      string
	EasyMonth           int
	FeedbackForChart    common.DifficultyChart
	Free                bool
	Id                  int
	Ip                  string
	IsTodo              bool
	IsCompetitive       bool
	Labels              []Label
	Name                string
	Os                  string
	PlayInfo            MachinePlayInfo
	Points              int
	Poweroff            int
	Recommended         int
	Release             time.Time
	RootOwnsCount       int
	SpFlag              int
	Star                float32
	StaticPoints        int
	UserOwnsCount       int
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

type MachinePlayInfo struct {
	ExpiresAt string
	IsActive  bool
}

type PlayInfo struct {
	ActivePlayerCount int
	ExpiresAt         time.Time
	IsActive          bool
	IsSpawned         bool
	IsSpawning        bool
}
type PlayInfoAlt struct {
	ExpiresAt time.Time
	Ip        string
	Ports     []int
	Status    string
}

type MachineProfileInfo struct {
	AcademyModules             []common.AcademyModule
	Active                     bool
	AuthUserFirstRootTime      string
	AuthUserFirstUserTime      string
	AuthUserHasReviewed        bool
	AuthUserHasSubmittedMatrix bool
	AuthUserInRootOwns         bool
	AuthUserInUserOwns         bool
	Avatar                     string
	CanAccessWalkthrough       bool
	DifficultyText             string
	FeedbackForChart           common.DifficultyChart
	Free                       bool
	HasChangelog               bool
	Id                         int
	InfoStatus                 string
	Ip                         string
	IsGuidedEnabled            bool
	IsTodo                     bool
	MachineMode                string
	Maker                      common.Maker
	Maker2                     common.Maker
	Name                       string
	Os                         string
	OwnRank                    int
	PlayInfo                   PlayInfo
	Points                     int
	Recommended                bool
	Release                    time.Time
	Retired                    bool
	ReviewsCount               int
	RootBlood                  common.BloodInfo
	RootOwnsCount              int
	SeasonId                   int
	ShowGoVip                  bool
	ShowGoVipServer            bool
	SpFlag                     int
	Stars                      float32
	StartMode                  string
	StaticPoints               int
	Synopsis                   string
	UserBlood                  common.BloodInfo
	UserCanReview              bool
	UserOwnsCount              int
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

type MachineRetiring struct {
	Avatar         string
	DifficultyText string
	Id             int
	Name           string
	Os             string
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

type MachineUnreleasedData struct {
	Avatar         string
	CoCreators     []common.UserBasicInfo
	Difficulty     int
	DifficultyText string
	FirstCreator   []common.UserBasicInfo
	Id             int
	Name           string
	Os             string
	Release        time.Time
	Retiring       MachineRetiring
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
