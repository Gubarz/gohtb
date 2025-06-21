package users

import (
	"context"

	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	"github.com/gubarz/gohtb/internal/service"
)

func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

func (s *Service) User(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

func (h *Handle) ProfileBasic(ctx context.Context) (ProfileBasicResponse, error) {
	resp, err := h.client.V4().GetUserProfileBasicWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) ProfileBasicResponse {
			return ProfileBasicResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return ProfileBasicResponse{
		Data: fromAPIUserProfile(resp.JSON200.Profile),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) ProfileActivity(ctx context.Context) (ProfileActivityResposnse, error) {
	resp, err := h.client.V4().GetUserProfileActivityWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) ProfileActivityResposnse {
			return ProfileActivityResposnse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return ProfileActivityResposnse{
		Data: convert.Slice(*resp.JSON200.Profile.Activity, fromAPIUserActivity),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}
