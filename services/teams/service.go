package teams

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

type Service struct {
	base service.Base
}

func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

type Handle struct {
	client service.Client
	id     int
}

// Team returns a handle for a specific team with the given ID.
// This handle can be used to perform operations related to that team,
// such as retrieving members, invitations, and activity data.
func (s *Service) Team(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

type UserEntry = v4Client.UserEntry

type InvitationsResponse struct {
	Data         []UserEntry
	ResponseMeta common.ResponseMeta
}

// Invitations retrieves pending invitations for the team.
// This returns a list of users who have been invited to join the team
// but have not yet accepted or rejected the invitation.
//
// Example:
//
//	invitations, err := client.Teams.Team(12345).Invitations(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, invite := range invitations.Data {
//		fmt.Printf("Pending invite: %s\n", invite.Username)
//	}
func (h *Handle) Invitations(ctx context.Context) (InvitationsResponse, error) {
	resp, err := h.client.V4().GetTeamInvitations(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	if err != nil {
		return InvitationsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetTeamInvitationsResponse)
	if err != nil {
		return InvitationsResponse{ResponseMeta: meta}, err
	}
	return InvitationsResponse{
		Data:         parsed.JSON200.Original,
		ResponseMeta: meta,
	}, nil
}

type TeamMember = v4Client.TeamMember

type MembersResponse struct {
	Data         []TeamMember
	ResponseMeta common.ResponseMeta
}

// Members retrieves the current members of the team.
// This returns a list of all users who are currently part of the team,
// including their roles and membership information.
//
// Example:
//
//	members, err := client.Teams.Team(12345).Members(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, member := range members.Data {
//		fmt.Printf("Member: %s (Role: %s)\n", member.Username, member.Role)
//	}
func (h *Handle) Members(ctx context.Context) (MembersResponse, error) {
	resp, err := h.client.V4().GetTeamMembers(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return MembersResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetTeamMembersResponse)
	if err != nil {
		return MembersResponse{ResponseMeta: meta}, err
	}

	return MembersResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type TeamActivityItem = v4Client.TeamActivityItem

type ActivityResponse struct {
	Data         []TeamActivityItem
	ResponseMeta common.ResponseMeta
}

type TeamStats = v4Client.TeamStatsOwnsTeamIdResponse
type StatsResponse struct {
	Data TeamStats
	ResponseMeta common.ResponseMeta
}

type TeamInfo = v4Client.TeamInfoTeamIdResponse

type TeamInfoResponse struct {
	Data TeamInfo
	ResponseMeta common.ResponseMeta
}

// Activity retrieves the activity history for the team.
// This includes recent team actions, achievements, and other team-related
// activities on the HackTheBox platform.
//
// Example:
//
//	activity, err := client.Teams.Team(12345).Activity(ctx, 30)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, act := range activity.Data {
//		fmt.Printf("Activity: %s at %s\n", act.Type, act.Date)
//	}
func (h *Handle) Activity(ctx context.Context, days int) (ActivityResponse, error) {
	// This is set to 90 days wich is max by the API
	// Max items returned is 100
	last := 90
	if days >= 1 && days <= 90 {
		last = days
	}
	params := &v4Client.GetTeamActivityParams{
		NPastDays: &last,
	}
	resp, err := h.client.V4().GetTeamActivity(
		h.client.Limiter().Wrap(ctx),
		h.id,
		params,
	)

	if err != nil {
		return ActivityResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetTeamActivityResponse)
	if err != nil {
		return ActivityResponse{ResponseMeta: meta}, err
	}
	return ActivityResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

// AcceptInvite accepts a team invitation with the specified ID.
// This operation adds the current user to the team that sent the invitation.
// This is the request ID not the User ID that sent the request.
//
// Example:
//
//	result, err := client.Teams.AcceptInvite(ctx, 67890)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Invite result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (s *Service) AcceptInvite(ctx context.Context, id int) (common.MessageResponse, error) {
	resp, err := s.base.Client.V4().PostTeamInviteAccept(
		s.base.Client.Limiter().Wrap(ctx),
		id,
	)

	if err != nil {
		return common.MessageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostTeamInviteAcceptResponse)
	if err != nil {
		return common.MessageResponse{ResponseMeta: meta}, err
	}

	return common.MessageResponse{
		Data: common.Message{
			Message: parsed.JSON200.Message,
			Success: parsed.JSON200.Success,
		},
		ResponseMeta: meta,
	}, nil
}

// RejectInvite rejects a team invitation with the specified ID.
// This operation declines the team invitation without joining the team.
// This is the request ID not the User ID that sent the request.
//
// Example:
//
//	result, err := client.Teams.RejectInvite(ctx, 67890)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Reject result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (s *Service) RejectInvite(ctx context.Context, id int) (common.MessageResponse, error) {
	resp, err := s.base.Client.V4().DeleteTeamInviteReject(
		s.base.Client.Limiter().Wrap(ctx),
		id,
	)
	if err != nil {
		return common.MessageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseDeleteTeamInviteRejectResponse)
	if err != nil {
		return common.MessageResponse{ResponseMeta: meta}, err
	}
	return common.MessageResponse{
		Data: common.Message{
			Message: parsed.JSON200.Message,
			Success: parsed.JSON200.Success,
		},
		ResponseMeta: meta,
	}, nil
}

// KickMember removes a user from the team with the specified user ID.
// This operation requires appropriate permissions within the team.
//
// Example:
//
//	result, err := client.Teams.KickMember(ctx, 54321)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Kick result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (s *Service) KickMember(ctx context.Context, id int) (common.MessageResponse, error) {
	resp, err := s.base.Client.V4().PostTeamKickUser(
		s.base.Client.Limiter().Wrap(ctx),
		id,
	)
	if err != nil {
		return common.MessageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostTeamKickUserResponse)
	if err != nil {
		return common.MessageResponse{ResponseMeta: meta}, err
	}

	return common.MessageResponse{
		Data: common.Message{
			Message: parsed.JSON200.Message,
			Success: parsed.JSON200.Success,
		},
		ResponseMeta: meta,
	}, nil
}


func (s *Service) Stats(ctx context.Context, id int) (StatsResponse, error) {
	resp, err := s.base.Client.V4().GetTeamStatsOwns(
		s.base.Client.Limiter().Wrap(ctx),
		id,
	)
	if err != nil {
		return StatsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetTeamStatsOwnsResponse)
	if err != nil {
		return StatsResponse{ResponseMeta: meta}, err
	}
	return StatsResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}


func (s *Service) TeamInfo(ctx context.Context, id int) (TeamInfoResponse, error) {
	resp, err := s.base.Client.V4().GetTeamInfo(
		s.base.Client.Limiter().Wrap(ctx),
		id,
	)
	if err != nil {
		return TeamInfoResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetTeamInfoResponse)
	if err != nil {
		return TeamInfoResponse{ResponseMeta: meta}, err
	}
	return TeamInfoResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}