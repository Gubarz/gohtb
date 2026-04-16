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

// NewService creates a new prolabs service bound to a shared client.
//
// Example:
//
//	prolabService := prolabs.NewService(client)
//	_ = prolabService
func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

type Handle struct {
	client service.Client
	id     int
}

// Prolab returns a handle for a specific prolab with the given ID.
//
// Example:
//
//	prolab := client.Prolabs.Prolab(1)
//	_ = prolab
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

// List retrieves the available prolabs.
//
// Example:
//
//	prolabs, err := client.Prolabs.List(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Prolabs found: %d\n", len(prolabs.Data))
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

// FAQ retrieves frequently asked questions for the selected prolab.
//
// Example:
//
//	faq, err := client.Prolabs.Prolab(1).FAQ(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("FAQ entries: %d\n", len(faq.Data))
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

// Flags retrieves available flags for the selected prolab.
//
// Example:
//
//	flags, err := client.Prolabs.Prolab(1).Flags(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Flags available: %d\n", len(flags.Data))
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

type ProlabData struct {
	v4Client.ProlabData
	DescriptionHTML string
}

func wrapProlabData(x v4Client.ProlabData) ProlabData {
	return ProlabData{ProlabData: x}
}

type InfoResponse struct {
	Data         ProlabData
	ResponseMeta common.ResponseMeta
}

// Info retrieves detailed information for the selected prolab.
//
// Example:
//
//	info, err := client.Prolabs.Prolab(1).Info(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Prolab: %s\n", info.Data.Name)
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
	wrapped := wrapProlabData(parsed.JSON200.Data)
	wrapped.DescriptionHTML = common.SanitizeHTML(wrapped.Description)
	wrapped.Description = common.StrictHTML(wrapped.Description)

	return InfoResponse{
		Data:         wrapped,
		ResponseMeta: meta,
	}, nil
}

type Machine = v4Client.Machine

type ProlabMachineData = []Machine

type MachinesResponse struct {
	Data         ProlabMachineData
	ResponseMeta common.ResponseMeta
}

// Machines retrieves machines linked to the selected prolab.
//
// Example:
//
//	machines, err := client.Prolabs.Prolab(1).Machines(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Prolab machines: %d\n", len(machines.Data))
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

// Overview retrieves summary statistics for the selected prolab.
//
// Example:
//
//	overview, err := client.Prolabs.Prolab(1).Overview(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Prolab users: %d\n", overview.Data.Users)
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

// Progress retrieves the authenticated user's progress for the selected prolab.
//
// Example:
//
//	progress, err := client.Prolabs.Prolab(1).Progress(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Progress percent: %d\n", progress.Data.CompletionPercentage)
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
	Data         float32
	ResponseMeta common.ResponseMeta
}

type MessageStatus struct {
	Message string
	Status  int
}

type SubmitFlagResponse struct {
	Data         MessageStatus
	ResponseMeta common.ResponseMeta
}

// SubmitFlag submits a flag for the selected prolab.
//
// Example:
//
//	result, err := client.Prolabs.Prolab(1).SubmitFlag(ctx, "HTB{example_flag}")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Submit result: %s\n", result.Data.Message)
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

// ChangelogsData contains prolab changelog entries.
type ChangelogsData struct {
	Data   []map[string]interface{}
	Status bool
}

// ChangelogsResponse contains prolab changelog payload.
type ChangelogsResponse struct {
	Data         ChangelogsData
	ResponseMeta common.ResponseMeta
}

// Changelogs retrieves changelog entries for the selected prolab.
//
// Example:
//
//	changelogs, err := client.Prolabs.Prolab(1).Changelogs(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Prolab changelog entries: %d\n", len(changelogs.Data.Data))
func (h *Handle) Changelogs(ctx context.Context) (ChangelogsResponse, error) {
	resp, err := h.client.V4().GetProlabChangelogs(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return ChangelogsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetProlabChangelogsResponse)
	if err != nil {
		return ChangelogsResponse{ResponseMeta: meta}, err
	}

	return ChangelogsResponse{
		Data: ChangelogsData{
			Data:   parsed.JSON200.Data,
			Status: parsed.JSON200.Status,
		},
		ResponseMeta: meta,
	}, nil
}
