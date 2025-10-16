package sherlocks

import (
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
	status     v4Client.GetSherlocksParamsStatus
	state      v4Client.State
	sortBy     v4Client.GetSherlocksParamsSortBy
	sortType   v4Client.GetSherlocksParamsSortType
	difficulty v4Client.Difficulty
	category   v4Client.Category
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

type DownloadResponse struct {
	Data         SherlockDownloadLink
	ResponseMeta common.ResponseMeta
}

type ProgressResponse struct {
	Data         SherlockProgressData
	ResponseMeta common.ResponseMeta
}

type TasksResponse struct {
	Data         SherlockTasksData
	ResponseMeta common.ResponseMeta
}

type SherlockTasksData = []SherlockTask

type OwnResponse struct {
	Data         TaskFlagResponse
	ResponseMeta common.ResponseMeta
}

type PlayResponse struct {
	Data         SherlockPlay
	ResponseMeta common.ResponseMeta
}

type SherlockItem = v4Client.SherlockItem
type SherlockRetires = v4Client.SherlockRetires
type SherlockNamedItemData = v4Client.SherlockNamedItemData
type SherlockDownloadLink = v4Client.SherlockDownloadLink
type SherlockProgressData = v4Client.SherlockProgressData
type SherlockTask = v4Client.SherlockTask
type TaskType = v4Client.TaskType
type TaskFlagResponse = v4Client.TaskFlagResponse
type UserTask = v4Client.UserTask
type SherlockPlay = v4Client.SherlockPlay
