package teams

import (
	"context"

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

func (s *Service) Team(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

func (h *Handle) Invitations(ctx context.Context) (InvitationsResponse, error) {
	resp, err := h.client.V4().GetTeamInvitationsWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, resp.StatusCode(), func(raw []byte) InvitationsResponse {
			return InvitationsResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return InvitationsResponse{
		Data: convert.Slice(*resp.JSON200.Original, fromAPIUserEntry),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Members(ctx context.Context) (MembersResponse, error) {
	resp, err := h.client.V4().GetTeamMembersWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, resp.StatusCode(), func(raw []byte) MembersResponse {
			return MembersResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return MembersResponse{
		Data: convert.Slice(*resp.JSON200, fromAPITeamMember),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Activity(ctx context.Context) (ActivityResponse, error) {
	resp, err := h.client.V4().GetTeamActivityWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, resp.StatusCode(), func(raw []byte) ActivityResponse {
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

func (s *Service) AcceptInvite(ctx context.Context, id int) (common.MessageResponse, error) {
	resp, err := s.base.Client.V4().PostTeamInviteAcceptWithResponse(
		s.base.Client.Limiter().Wrap(ctx),
		id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, resp.StatusCode(), func(raw []byte) common.MessageResponse {
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

func (s *Service) RejectInvite(ctx context.Context, id int) (common.MessageResponse, error) {
	resp, err := s.base.Client.V4().DeleteTeamInviteRejectWithResponse(
		s.base.Client.Limiter().Wrap(ctx),
		id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, resp.StatusCode(), func(raw []byte) common.MessageResponse {
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

func (s *Service) KickMember(ctx context.Context, id int) (common.MessageResponse, error) {
	resp, err := s.base.Client.V4().PostTeamKickUserWithResponse(
		s.base.Client.Limiter().Wrap(ctx),
		id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, resp.StatusCode(), func(raw []byte) common.MessageResponse {
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
