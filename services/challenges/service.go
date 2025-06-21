package challenges

import (
	"context"
	"strconv"

	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/deref"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	v4client "github.com/gubarz/gohtb/internal/httpclient/v4"
	"github.com/gubarz/gohtb/internal/service"
)

func NewService(client service.Client, product string) *Service {
	return &Service{
		base:    service.NewBase(client),
		product: product,
	}
}

func (s *Service) Challenge(id int) *Handle {
	return &Handle{
		client:  s.base.Client,
		id:      id,
		product: s.product,
	}
}

func (s *Service) List() *ChallengeQuery {
	return &ChallengeQuery{
		client:  s.base.Client,
		page:    1,
		perPage: 100,
	}
}

func (h *Handle) Info(ctx context.Context) (InfoResponse, error) {
	slug := strconv.Itoa(h.id)
	resp, err := h.client.V4().GetChallengeInfoWithResponse(
		h.client.Limiter().Wrap(ctx),
		slug,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil || resp.JSON200.Challenge == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) InfoResponse {
			return InfoResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}
	return InfoResponse{
		Data: fromAPIChallengeInfo(*resp.JSON200.Challenge),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) ToDo(ctx context.Context) (common.TodoUpdateResponse, error) {
	resp, err := h.client.V4().PostTodoUpdateWithResponse(
		h.client.Limiter().Wrap(ctx),
		v4client.PostTodoUpdateParamsProduct(h.product),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil || resp.JSON200.Info == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) common.TodoUpdateResponse {
			return common.TodoUpdateResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}
	return common.TodoUpdateResponse{
		Data: convert.Slice(*resp.JSON200.Info, common.FromAPIInfoArray),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Start(ctx context.Context) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostChallengeStartWithFormdataBodyWithResponse(
		h.client.Limiter().Wrap(ctx),
		v4client.PostChallengeStartFormdataRequestBody{
			ChallengeId: h.id,
		},
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
		},
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Stop(ctx context.Context) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostChallengeStopWithFormdataBodyWithResponse(
		h.client.Limiter().Wrap(ctx),
		v4client.PostChallengeStopFormdataRequestBody{
			ChallengeId: h.id,
		},
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil || resp.JSON200.Message == nil {
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

func (h *Handle) Own(ctx context.Context, flag string) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostChallengeOwnWithFormdataBodyWithResponse(
		h.client.Limiter().Wrap(ctx),
		v4client.ChallengeOwnRequest{
			ChallengeId: h.id,
			Flag:        flag,
		},
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil || resp.JSON200.Message == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) common.MessageResponse {
			return common.MessageResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}
	return common.MessageResponse{
		Data: common.Message{
			Message: deref.String(resp.JSON200.Message),
		},
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Activity(ctx context.Context) (ActivityResponse, error) {
	resp, err := h.client.V4().GetChallengeActivityWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil || resp.JSON200.Info == nil || resp.JSON200.Info.Activity == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) ActivityResponse {
			return ActivityResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}
	return ActivityResponse{
		Data: convert.Slice(*resp.JSON200.Info.Activity, fromAPIChallengeActivity),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Download(ctx context.Context) (DownloadResponse, error) {
	resp, err := h.client.V4().GetChallengeDownloadWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) DownloadResponse {
			return DownloadResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}
	return DownloadResponse{
		Data: raw,
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}
