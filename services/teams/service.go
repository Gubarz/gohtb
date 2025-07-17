package teams

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/deref"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	"github.com/gubarz/gohtb/internal/service"
)

func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
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
	resp, err := h.client.V4().GetTeamInvitationsWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) InvitationsResponse {
			return InvitationsResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return InvitationsResponse{
		Data: convert.SlicePointer(resp.JSON200.Original, fromAPIUserEntry),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
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
	resp, err := h.client.V4().GetTeamMembersWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) MembersResponse {
			return MembersResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return MembersResponse{
		Data: convert.SlicePointer(resp.JSON200, fromAPITeamMember),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

// Activity retrieves the activity history for the team.
// This includes recent team actions, achievements, and other team-related
// activities on the HackTheBox platform.
//
// Example:
//
//	activity, err := client.Teams.Team(12345).Activity(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, act := range activity.Data {
//		fmt.Printf("Activity: %s at %s\n", act.Type, act.Date)
//	}
func (h *Handle) Activity(ctx context.Context) (ActivityResponse, error) {
	// This is set to 90 days wich is max by the API
	// Not setting it can pausibly return data that is not up to date
	// Max items returned is 100
	last := 90
	params := &v4Client.GetTeamActivityParams{
		NPastDays: &last,
	}
	resp, err := h.client.V4().GetTeamActivityWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
		params,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) ActivityResponse {
			return ActivityResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return ActivityResponse{
		Data: convert.Slice(*resp.JSON200, fromAPITeamActivityItem),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
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
	resp, err := s.base.Client.V4().PostTeamInviteAcceptWithResponse(
		s.base.Client.Limiter().Wrap(ctx),
		id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) common.MessageResponse {
			return common.MessageResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return common.MessageResponse{
		Data: common.Message{
			Message: deref.String(resp.JSON200.Message),
			Success: deref.Bool(resp.JSON200.Success),
		},
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
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
	resp, err := s.base.Client.V4().DeleteTeamInviteRejectWithResponse(
		s.base.Client.Limiter().Wrap(ctx),
		id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) common.MessageResponse {
			return common.MessageResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return common.MessageResponse{
		Data: common.Message{
			Message: deref.String(resp.JSON200.Message),
			Success: deref.Bool(resp.JSON200.Success),
		},
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
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
	resp, err := s.base.Client.V4().PostTeamKickUserWithResponse(
		s.base.Client.Limiter().Wrap(ctx),
		id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) common.MessageResponse {
			return common.MessageResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return common.MessageResponse{
		Data: common.Message{
			Message: deref.String(resp.JSON200.Message),
			Success: deref.Bool(resp.JSON200.Success),
		},
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}
