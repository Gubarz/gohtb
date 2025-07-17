package prolabs

import (
	"time"

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

type FaqItem struct {
	Answer   string
	Generic  bool
	Question string
}

type ProlabsData struct {
	Count int
	Labs  ProlabDataItems
}

type ProlabDataItems = []Prolab

type Prolab struct {
	CoverImgUrl                string
	DesignatedCategory         string
	Id                         int
	Identifier                 string
	LabServersCount            int
	Level                      int
	Mini                       bool
	Name                       string
	New                        bool
	Ownership                  float32
	ProFlagsCount              int
	ProMachinesCount           int
	ReleaseAt                  time.Time
	SkillLevel                 string
	State                      string
	Team                       string
	UserEligibleForCertificate bool
}

type FlagsItems = []common.Flag

type ProlabData struct {
	ActiveUsers      int
	CanInteract      bool
	CoverImageUrl    string
	Description      string
	EntryPoints      common.StringArray
	Id               int
	Identifier       string
	LabMasters       LabMasterItems
	LabServersCount  int
	Mini             bool
	Name             string
	ProFlagsCount    int
	ProMachinesCount int
	State            string
	Version          string
	VideoUrl         string
	Writeup          string
}

type LabMasterItems = []common.UserIdNameThumb

type ProlabMachineData = []Machine

type Machine struct {
	AvatarThumbUrl string
	Id             int
	Name           string
	Os             string
}

type ProlabOverviewData struct {
	DesignatedLevel    DesignatedLevel
	Excerpt            string
	Id                 int
	Identifier         string
	LabMasters         LabMasterItems
	Mini               bool
	Name               string
	NewVersion         bool
	OverviewImageUrl   string
	ProFlagsCount      int
	ProMachinesCount   int
	SkillLevel         string
	SocialLinks        SocialLinks
	State              string
	UserEligibleToPlay bool
	Version            string
}

type DesignatedLevel struct {
	Category    string
	Description string
	Level       int
	Team        string
}

type SocialLinks struct {
	Discord string
	Forum   string
}

type ProlabProgressData struct {
	KeyedProLabMileStone              KeyedProlabMileStoneItems
	Ownership                         float32
	OwnershipRequiredForCertification float32
	UserEligibleForCertificate        bool
}

type KeyedProlabMileStoneItems = []KeyedProLabMileStone

type KeyedProLabMileStone struct {
	Description        string
	Icon               string
	IsMilestoneReached bool
	Percent            float32
	Rarity             float32
	Text               string
}

type ProlabSubscription struct {
	Active             bool
	EndsAt             string
	Name               string
	RenewsAt           string
	SubscriptionPeriod string
	Type               string
}
