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

// NewService creates a new teams service bound to a shared client.
//
// Example:
//
//	teamService := teams.NewService(client)
//	_ = teamService
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
//
// Example:
//
//	team := client.Teams.Team(12345)
//	_ = team
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

type TeamMember struct {
	v4Client.TeamMember
	Rank int
}

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
		Data:         wrapMembers(*parsed.JSON200),
		ResponseMeta: meta,
	}, nil
}

func wrapMembers(list []v4Client.TeamMember) []TeamMember {
	out := make([]TeamMember, len(list))
	for i, m := range list {
		out[i] = TeamMember{
			TeamMember: m,
			Rank:       rankToInt(m.Rank),
		}
	}
	return out
}

func rankToInt(u v4Client.TeamMember_Rank) int {
	n, err := u.AsTeamMemberRank0()
	if err != nil {
		return 0
	}
	return n
}

type TeamActivityItem = v4Client.TeamActivityItem

type ActivityResponse struct {
	Data         []TeamActivityItem
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
	// This is set to 90 days which is max by the API
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

type TeamStats = v4Client.TeamStatsOwnsTeamIdResponse
type TeamStatsResponse struct {
	Data         TeamStats
	ResponseMeta common.ResponseMeta
}

// Stats gets the team stats like user owns, system owns, bloods, points
//
// Example:
//
//	stats, err := client.Teams.Team(12345).Stats(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Team stats: bloods=%d systemOwns=%d userOwns=%d\n", stats.Data.FirstBloods, stats.Data.SystemOwns, stats.Data.UserOwns)
func (h *Handle) Stats(ctx context.Context) (TeamStatsResponse, error) {
	resp, err := h.client.V4().GetTeamStatsOwns(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return TeamStatsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}
	parsed, meta, err := common.Parse(resp, v4Client.ParseGetTeamStatsOwnsResponse)
	if err != nil {
		return TeamStatsResponse{ResponseMeta: meta}, err
	}
	return TeamStatsResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type TeamInfo = v4Client.TeamInfoTeamIdResponse

type TeamInfoResponse struct {
	Data         TeamInfo
	ResponseMeta common.ResponseMeta
}

// Team info gets the team information like members, motto, description, country, etc.
//
// Example:
//
//	info, err := client.Teams.Team(12345).Info(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Team: %s (Motto: %s)\n", info.Data.Name, info.Data.Motto)
func (h *Handle) Info(ctx context.Context) (TeamInfoResponse, error) {
	resp, err := h.client.V4().GetTeamInfo(
		h.client.Limiter().Wrap(ctx),
		h.id,
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

type ChartChallengeCategoriesData = v4Client.TeamChartChallengeTeamIdResponse

// ChartChallengeCategoriesResponse contains challenge category chart data.
type ChartChallengeCategoriesResponse struct {
	Data         ChartChallengeCategoriesData
	ResponseMeta common.ResponseMeta
}

// ChartChallengeCategories retrieves challenge category chart data for the team.
//
// Example:
//
//	chart, err := client.Teams.Team(12345).ChartChallengeCategories(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Challenge category chart: %+v\n", chart.Data)
func (h *Handle) ChartChallengeCategories(ctx context.Context) (ChartChallengeCategoriesResponse, error) {
	resp, err := h.client.V4().GetTeamChartChallengeCategories(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return ChartChallengeCategoriesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetTeamChartChallengeCategoriesResponse)
	if err != nil {
		return ChartChallengeCategoriesResponse{ResponseMeta: meta}, err
	}

	return ChartChallengeCategoriesResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ChartMachinesAttackData = v4Client.TeamChartMachinesAttackTeamIdResponse

// ChartMachinesAttackResponse contains machine attack chart data.
type ChartMachinesAttackResponse struct {
	Data         ChartMachinesAttackData
	ResponseMeta common.ResponseMeta
}

// ChartMachinesAttack retrieves machine attack chart data for the team.
//
// Example:
//
//	chart, err := client.Teams.Team(12345).ChartMachinesAttack(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Machine attack chart: %+v\n", chart.Data)
func (h *Handle) ChartMachinesAttack(ctx context.Context) (ChartMachinesAttackResponse, error) {
	resp, err := h.client.V4().GetTeamChartMachinesAttack(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return ChartMachinesAttackResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetTeamChartMachinesAttackResponse)
	if err != nil {
		return ChartMachinesAttackResponse{ResponseMeta: meta}, err
	}

	return ChartMachinesAttackResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type GraphData = v4Client.TeamGraphTeamIdResponse

// GraphResponse contains team graph data.
type GraphResponse struct {
	Data         GraphData
	ResponseMeta common.ResponseMeta
}

// Graph retrieves graph data for the team.
//
// Example:
//
//	graph, err := client.Teams.Team(12345).Graph(ctx, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Team graph: %+v\n", graph.Data)
func (h *Handle) Graph(ctx context.Context, params *v4Client.GetTeamGraphParams) (GraphResponse, error) {
	resp, err := h.client.V4().GetTeamGraph(h.client.Limiter().Wrap(ctx), h.id, params)
	if err != nil {
		return GraphResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetTeamGraphResponse)
	if err != nil {
		return GraphResponse{ResponseMeta: meta}, err
	}

	return GraphResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}
