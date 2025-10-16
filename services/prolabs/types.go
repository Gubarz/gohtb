package prolabs

import (
	v4client "github.com/gubarz/gohtb/httpclient/v4"
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

type ListResponse struct {
	Data         ProlabsData
	ResponseMeta common.ResponseMeta
}

type FaqResponse struct {
	Data         ProlabFaqData
	ResponseMeta common.ResponseMeta
}

type MachinesResponse struct {
	Data         ProlabMachineData
	ResponseMeta common.ResponseMeta
}

type OverviewResponse struct {
	Data         ProlabOverviewData
	ResponseMeta common.ResponseMeta
}

type ProgressResponse struct {
	Data         ProlabProgressData
	ResponseMeta common.ResponseMeta
}

type RatingResponse struct {
	Data         string
	ResponseMeta common.ResponseMeta
}

type SubscriptionResponse struct {
	Data         ProlabSubscription
	ResponseMeta common.ResponseMeta
}

type SubmitFlagResponse struct {
	Data         MessageStatus
	ResponseMeta common.ResponseMeta
}

type MessageStatus struct {
	Message string
	Status  int
}

type FlagsResponse struct {
	Data         []common.Flag
	ResponseMeta common.ResponseMeta
}

type InfoResponse struct {
	Data         ProlabData
	ResponseMeta common.ResponseMeta
}

type ProlabFaqData = []FaqItem

type ProlabDataItems = []Prolab

type FlagsItems = []common.Flag

type LabMasterItems = []common.UserIdNameThumb

type ProlabMachineData = []Machine

type KeyedProlabMileStoneItems = []KeyedProLabMileStone

type Prolab = v4client.Prolab
type FaqItem = v4client.FaqItem
type ProlabData = v4client.ProlabData
type ProlabOverviewData = v4client.ProlabOverviewData
type Machine = v4client.Machine
type DesignatedLevel = v4client.DesignatedLevel
type SocialLinks = v4client.SocialLinks
type ProlabProgressData = v4client.ProlabProgressData
type KeyedProLabMileStone = v4client.KeyedProLabMileStone
type ProlabSubscription = v4client.ProlabSubscription
type ProlabsData = v4client.ProlabsData
