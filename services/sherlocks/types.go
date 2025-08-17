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
	Retires             time.Time
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
