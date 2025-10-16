package prolabs

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

func (s *Service) Prolab(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

type ProlabsData = v4Client.ProlabsData

type ListResponse struct {
	Data         ProlabsData
	ResponseMeta common.ResponseMeta
}

func (s *Service) List(ctx context.Context) (ListResponse, error) {
	resp, err := s.base.Client.V4().GetProlabs(
		s.base.Client.Limiter().Wrap(ctx))

	if err != nil {
		return ListResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetProlabsResponse)
	if err != nil {
		return ListResponse{ResponseMeta: meta}, err
	}

	return ListResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type FaqItem = v4Client.FaqItem

type ProlabFaqData = []FaqItem

type FaqResponse struct {
	Data         ProlabFaqData
	ResponseMeta common.ResponseMeta
}

func (h *Handle) FAQ(ctx context.Context) (FaqResponse, error) {
	resp, err := h.client.V4().GetProlabFaq(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return FaqResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetProlabFaqResponse)
	if err != nil {
		return FaqResponse{ResponseMeta: meta}, err
	}

	return FaqResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type FlagsResponse struct {
	Data         []common.Flag
	ResponseMeta common.ResponseMeta
}

func (h *Handle) Flags(ctx context.Context) (FlagsResponse, error) {
	resp, err := h.client.V4().GetProlabFlags(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return FlagsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetProlabFlagsResponse)
	if err != nil {
		return FlagsResponse{ResponseMeta: meta}, err
	}

	return FlagsResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type ProlabData = v4Client.ProlabData

type InfoResponse struct {
	Data         ProlabData
	ResponseMeta common.ResponseMeta
}

func (h *Handle) Info(ctx context.Context) (InfoResponse, error) {
	resp, err := h.client.V4().GetProlabInfo(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return InfoResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetProlabInfoResponse)
	if err != nil {
		return InfoResponse{ResponseMeta: meta}, err
	}

	return InfoResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type Machine = v4Client.Machine

type ProlabMachineData = []Machine

type MachinesResponse struct {
	Data         ProlabMachineData
	ResponseMeta common.ResponseMeta
}

func (h *Handle) Machines(ctx context.Context) (MachinesResponse, error) {
	resp, err := h.client.V4().GetProlabMachines(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return MachinesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetProlabMachinesResponse)
	if err != nil {
		return MachinesResponse{ResponseMeta: meta}, err
	}

	return MachinesResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type ProlabOverviewData = v4Client.ProlabOverviewData

type OverviewResponse struct {
	Data         ProlabOverviewData
	ResponseMeta common.ResponseMeta
}

func (h *Handle) Overview(ctx context.Context) (OverviewResponse, error) {
	resp, err := h.client.V4().GetProlabOverview(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return OverviewResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetProlabOverviewResponse)
	if err != nil {
		return OverviewResponse{ResponseMeta: meta}, err
	}
	return OverviewResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type ProlabProgressData = v4Client.ProlabProgressData

type ProgressResponse struct {
	Data         ProlabProgressData
	ResponseMeta common.ResponseMeta
}

func (h *Handle) Progress(ctx context.Context) (ProgressResponse, error) {
	resp, err := h.client.V4().GetProlabProgress(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return ProgressResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetProlabProgressResponse)
	if err != nil {
		return ProgressResponse{ResponseMeta: meta}, err
	}

	return ProgressResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type RatingResponse struct {
	Data         string
	ResponseMeta common.ResponseMeta
}

func (h *Handle) Rating(ctx context.Context) (RatingResponse, error) {
	resp, err := h.client.V4().GetProlabRating(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return RatingResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetProlabRatingResponse)
	if err != nil {
		return RatingResponse{ResponseMeta: meta}, err
	}
	return RatingResponse{
		Data:         parsed.JSON200.Data.Rating,
		ResponseMeta: meta,
	}, nil
}

type ProlabSubscription = v4Client.ProlabSubscription

type SubscriptionResponse struct {
	Data         ProlabSubscription
	ResponseMeta common.ResponseMeta
}

func (h *Handle) Subscription(ctx context.Context) (SubscriptionResponse, error) {
	resp, err := h.client.V4().GetProlabSubscription(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return SubscriptionResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetProlabSubscriptionResponse)
	if err != nil {
		return SubscriptionResponse{ResponseMeta: meta}, err
	}
	return SubscriptionResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type MessageStatus struct {
	Message string
	Status  int
}

type SubmitFlagResponse struct {
	Data         MessageStatus
	ResponseMeta common.ResponseMeta
}

func (h *Handle) SubmitFlag(ctx context.Context, flag string) (SubmitFlagResponse, error) {
	resp, err := h.client.V4().PostProlabFlag(
		h.client.Limiter().Wrap(ctx),
		h.id,
		v4Client.PostProlabFlagJSONRequestBody{
			Flag: flag,
		},
	)
	if err != nil {
		return SubmitFlagResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostProlabFlagResponse)
	if err != nil {
		return SubmitFlagResponse{ResponseMeta: meta}, err
	}
	return SubmitFlagResponse{
		Data: MessageStatus{
			Message: parsed.JSON200.Message,
			Status:  parsed.JSON200.Status,
		},
		ResponseMeta: meta,
	}, nil
}
