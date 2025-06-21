package vms

import (
	"context"

	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/deref"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	v4client "github.com/gubarz/gohtb/internal/httpclient/v4"
	"github.com/gubarz/gohtb/internal/service"
)

func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

func (s *Service) VM(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

func (h *Handle) Reset(ctx context.Context) (Response, error) {
	resp, err := h.client.V4().PostVMResetWithResponse(
		h.client.Limiter().Wrap(ctx),
		v4client.PostVMResetJSONRequestBody{
			MachineId: h.id,
		})

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) Response {
			return Response{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return Response{
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

func (h *Handle) Spawn(ctx context.Context) (Response, error) {
	resp, err := h.client.V4().PostVMSpawnWithResponse(
		h.client.Limiter().Wrap(ctx),
		v4client.PostVMSpawnJSONRequestBody{
			MachineId: h.id,
		})

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) Response {
			return Response{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return Response{
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

func (h *Handle) Extend(ctx context.Context) (Response, error) {
	resp, err := h.client.V4().PostVMExtendWithResponse(
		h.client.Limiter().Wrap(ctx),
		v4client.PostVMExtendJSONRequestBody{
			MachineId: h.id,
		})

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) Response {
			return Response{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return Response{
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

func (h *Handle) Terminate(ctx context.Context) (Response, error) {
	resp, err := h.client.V4().PostVMTerminateWithResponse(
		h.client.Limiter().Wrap(ctx),
		v4client.PostVMTerminateJSONRequestBody{
			MachineId: h.id,
		})

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) Response {
			return Response{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return Response{
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

func (h *Handle) VoteReset(ctx context.Context) (Response, error) {
	resp, err := h.client.V4().PostVMResetVoteWithResponse(
		h.client.Limiter().Wrap(ctx),
		v4client.PostVMResetVoteJSONRequestBody{
			MachineId: h.id,
		})

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) Response {
			return Response{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return Response{
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

func (h *Handle) VoteResetAccept(ctx context.Context) (Response, error) {
	resp, err := h.client.V4().PostVMResetVoteAcceptWithResponse(
		h.client.Limiter().Wrap(ctx),
		v4client.PostVMResetVoteAcceptJSONRequestBody{
			MachineId: h.id,
		})

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) Response {
			return Response{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return Response{
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
