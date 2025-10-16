package challenges

import (
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
	status     v4Client.GetChallengesParamsStatus
	state      v4Client.State
	sortBy     v4Client.GetChallengesParamsSortBy
	sortType   v4Client.GetChallengesParamsSortType
	difficulty v4Client.Difficulty
	category   v4Client.Category
	todo       v4Client.GetChallengesParamsTodo
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

type ChallengeList = v4Client.ChallengeList
type ChallengeRetires = v4Client.ChallengeRetires
type ChallengeActivity = v4Client.ChallengeActivity
type Challenge = v4Client.Challenge
