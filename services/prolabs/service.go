package prolabs

import (
	"context"

	v4client "github.com/gubarz/gohtb/httpclient/v4"
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

func (s *Service) Prolab(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

func (s *Service) List(ctx context.Context) (ListResponse, error) {
	resp, err := s.base.Client.V4().GetProlabsWithResponse(
		s.base.Client.Limiter().Wrap(ctx))

	raw := extract.Raw(resp)
	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) ListResponse {
			return ListResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return ListResponse{
		Data: ProlabsData{
			Count: deref.Int(resp.JSON200.Data.Count),
			Labs:  convert.SlicePointer(resp.JSON200.Data.Labs, fromAPIProlab),
		},
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) FAQ(ctx context.Context) (FaqResponse, error) {
	resp, err := h.client.V4().GetProlabFaqWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) FaqResponse {
			return FaqResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return FaqResponse{
		Data: convert.SlicePointer(resp.JSON200.Data, fromAPIFaqItem),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Flags(ctx context.Context) (FlagsResponse, error) {
	resp, err := h.client.V4().GetProlabFlagsWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) FlagsResponse {
			return FlagsResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return FlagsResponse{
		Data: convert.SlicePointer(resp.JSON200.Data, common.FromAPIFlag),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Info(ctx context.Context) (InfoResponse, error) {
	resp, err := h.client.V4().GetProlabInfoWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) InfoResponse {
			return InfoResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return InfoResponse{
		Data: fromAPIProlabData(resp.JSON200.Data),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Machines(ctx context.Context) (MachinesResponse, error) {
	resp, err := h.client.V4().GetProlabMachinesWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) MachinesResponse {
			return MachinesResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return MachinesResponse{
		Data: convert.SlicePointer(resp.JSON200.Data, fromAPIProlabMachineData),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Overview(ctx context.Context) (OverviewResponse, error) {
	resp, err := h.client.V4().GetProlabOverviewWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) OverviewResponse {
			return OverviewResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return OverviewResponse{
		Data: fromAPIProlabOverviewData(resp.JSON200.Data),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Progress(ctx context.Context) (ProgressResponse, error) {
	resp, err := h.client.V4().GetProlabProgressWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) ProgressResponse {
			return ProgressResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return ProgressResponse{
		Data: fromAPIProlabProgressData(resp.JSON200.Data),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Rating(ctx context.Context) (RatingResponse, error) {
	resp, err := h.client.V4().GetProlabRatingWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) RatingResponse {
			return RatingResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return RatingResponse{
		Data: deref.String(resp.JSON200.Data.Rating),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Subscription(ctx context.Context) (SubscriptionResponse, error) {
	resp, err := h.client.V4().GetProlabSubscriptionWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) SubscriptionResponse {
			return SubscriptionResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return SubscriptionResponse{
		Data: fromAPIProlabSubscription(resp.JSON200.Data),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) SubmitFlag(ctx context.Context, flag string) (SubmitFlagResponse, error) {
	resp, err := h.client.V4().PostProlabFlagWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
		v4client.PostProlabFlagJSONRequestBody{
			Flag: flag,
		},
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) SubmitFlagResponse {
			return SubmitFlagResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return SubmitFlagResponse{
		Data: MessageStatus{
			Message: deref.String(resp.JSON200.Message),
			Status:  deref.Int(resp.JSON200.Status),
		},
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}
