package common

import (
	"net/http"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
)

type TodoItem = v4Client.Item
type Flag = v4Client.Flag
type UserBasicInfo = v4Client.UserBasicInfo
type Links = v4Client.Links
type Meta = v4Client.Meta
type MetaAlt = v4Client.MetaAlt
type PaginationLink = v4Client.PaginationLink
type Tag = v4Client.Tag
type TagCategory = v4Client.TagCategory
type HelpfulReviews = v4Client.HelpfulReviews
type DifficultyChart = v4Client.DifficultyChart
type Item = v4Client.Item
type PlayInfoAlt = v4Client.PlayInfoAlt
type BloodInfo = v4Client.BloodInfo
type Maker = v4Client.Maker
type MatrixInfo = v4Client.MatrixInfo
type AcademyTiers = v4Client.AcademyTiers
type TeamMachineAttackPaths = v4Client.TeamMachineAttackPaths
type TeamsAttackPathCard = v4Client.TeamsAttackPathCard
type UserIdNameThumb = v4Client.UserIdNameThumb
type AcademyDifficulty = v4Client.AcademyDifficulty
type AcademyModule = v4Client.AcademyModule
type UserBasicInfoWithRespect = v4Client.UserBasicInfoWithRespect
type PlayInfo = v4Client.PlayInfo

type MessageResponse struct {
	Data Message
	ResponseMeta
}

type ResponseMeta struct {
	Raw        []byte
	StatusCode int
	Headers    http.Header
	CFRay      string
}

type FlagData struct {
	Flags  []Flag
	Status bool
	ResponseMeta
}
type InfoArray = []Item

type TodoUpdateResponse struct {
	Data InfoArray
	ResponseMeta
}

type Message struct {
	Message string
	Success bool
}

type OwnsResponse struct {
	Data string
	ResponseMeta
}

type Messagesuccess struct {
	Message string
	Success bool
}

type IntArray = []int
type StringArray = []string
