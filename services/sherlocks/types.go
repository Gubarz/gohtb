package sherlocks

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

type SherlockQuery struct {
	client     service.Client
	status     *v4Client.GetSherlocksParamsStatus
	state      *v4Client.State
	sortBy     *v4Client.GetSherlocksParamsSortBy
	sortType   *v4Client.GetSherlocksParamsSortType
	difficulty *v4Client.Difficulty
	category   *v4Client.Category
	page       int
	perPage    int
}

type InfoResponse struct {
	Data         SherlockNamedItemData
	ResponseMeta common.ResponseMeta
}

type SherlockListResponse struct {
	Data         []SherlockItem
	ResponseMeta common.ResponseMeta
}

type SherlockListData = []SherlockItem

type SherlockItem struct {
	AuthUserHasReviewed bool
	Avatar              string
	CategoryId          int
	CategoryName        string
	Difficulty          string
	Id                  int
	IsOwned             bool
	Name                string
	Pinned              bool
	PlayMethods         []string
	Progress            int
	Rating              float32
	RatingCount         int
	ReleaseDate         string
	Retires             SherlockRetires
	Solves              int
	State               string
}

type SherlockNamedItemData struct {
	AuthUserHasReviewed bool
	Avatar              string
	CategoryId          int
	CategoryName        string
	Difficulty          string
	Favorite            bool
	Id                  int
	IsTodo              bool
	Name                string
	PlayMethods         []string
	Rating              float32
	RatingCount         int
	ReleaseAt           time.Time
	Retired             bool
	ShowGoVip           bool
	State               string
	Tags                []string
	UserCanReview       bool
	WriteupVisible      bool
}

type DownloadResponse struct {
	Data         SherlockDownloadLink
	ResponseMeta common.ResponseMeta
}

type SherlockDownloadLink struct {
	ExpiresIn int
	Url       string
}

type ProgressResponse struct {
	Data         SherlockProgressData
	ResponseMeta common.ResponseMeta
}

type SherlockProgressData struct {
	IsOwned       bool
	OwnRank       int
	Progress      int
	TasksAnswered int
	TotalTasks    int
}

type TasksResponse struct {
	Data         SherlockTasksData
	ResponseMeta common.ResponseMeta
}

type SherlockTasksData = []SherlockTask

type SherlockTask struct {
	Completed      bool
	Description    string
	Flag           string
	Hint           string
	Id             int
	MaskedFlag     string
	PrerequisiteId int
	TaskType       TaskType
	Title          string
	Type           TaskType
}

type TaskType struct {
	Id   int
	Text string
}

type OwnResponse struct {
	Data         TaskFlagResponse
	ResponseMeta common.ResponseMeta
}

type TaskFlagResponse struct {
	Message  string
	UserTask UserTask
}

type UserTask struct {
	Id     int
	TaskId int
	UserId int
}

type PlayResponse struct {
	Data         SherlockPlay
	ResponseMeta common.ResponseMeta
}

type SherlockPlay struct {
	Creators []common.Maker
	FileName string
	FileSize string
	Id       int
	PlayInfo common.PlayInfoAlt
	Scenario string
}

type SherlockRetires struct {
	AvatarUrl  string
	Difficulty string
	Name       string
}
