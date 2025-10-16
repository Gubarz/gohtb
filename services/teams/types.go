package teams

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

type InvitationsResponse struct {
	Data         []UserEntry
	ResponseMeta common.ResponseMeta
}

type MembersResponse struct {
	Data         []TeamMember
	ResponseMeta common.ResponseMeta
}

type ActivityResponse struct {
	Data         []TeamActivityItem
	ResponseMeta common.ResponseMeta
}

type UserEntry = v4Client.UserEntry
type User = v4Client.User
type UserRanking = v4Client.UserRanking
type TeamMember = v4Client.TeamMember
type TeamMemberTeam = v4Client.TeamMemberTeam
type TeamActivityItem = v4Client.TeamActivityItem
type TeamActivityUser = v4Client.TeamActivityUser
