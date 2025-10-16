package fortresses

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

type FlagData = common.FlagData
type Flag = common.Flag

type Info struct {
	Data         Data
	ResponseMeta common.ResponseMeta
}

type ListResponse struct {
	Data         []Item
	ResponseMeta common.ResponseMeta
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

type Data = v4Client.FortressData
type Company = v4Client.Company
type UserAvailability = v4Client.UserAvailability
type Item = v4Client.Fortress
