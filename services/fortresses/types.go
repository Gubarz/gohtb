package fortresses

import (
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

type FlagData = common.FlagData
type Flag = common.Flag

type Company struct {
	Description string
	Id          int
	Image       string
	Name        string
	Url         string
}

type UserAvailability struct {
	Available bool
	Code      int
	Message   string
}

type Data struct {
	Company              Company
	CompletionMessage    string
	CoverImageUrl        string
	Description          string
	Flags                []Flag
	HasCompletionMessage bool
	Id                   int
	Image                string
	Ip                   string
	Name                 string
	PlayersCompleted     int
	Points               int
	ProgressPercent      float32
	ResetVotes           int
	UserAvailability     UserAvailability
}

type Info struct {
	Data         Data
	ResponseMeta common.ResponseMeta
}

type ListResponse struct {
	Data         []Item
	ResponseMeta common.ResponseMeta
}

type Item struct {
	CoverImageUrl string
	Id            int
	Image         string
	Name          string
	New           bool
	NumberOfFlags int
}

type SubmitFlagResponse struct {
	Data         SubmitFlagData
	ResponseMeta common.ResponseMeta
}

type SubmitFlagData struct {
	Message string
	Status  int
}

type ResetResponse struct {
	Data         ResetFlagData
	ResponseMeta common.ResponseMeta
}

type ResetFlagData struct {
	Message string
	Status  bool
}
