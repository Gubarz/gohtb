package fortresses

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

type Service struct {
	base service.Base
}

// NewService creates a new fortresses service bound to a shared client.
//
// Example:
//
//	fortressService := fortresses.NewService(client)
//	_ = fortressService
func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

type Fortress = v4Client.Fortress

type ListResponse struct {
	Data         []Fortress
	ResponseMeta common.ResponseMeta
}

// List retrieves basic information about the fortresses.
//
// It returns both the normalized Data and the raw API response.
//
// Example:
//
//	fortresses, err := client.Fortresses.List(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Fortresses found: %d\n", len(fortresses.Data))
func (s *Service) List(ctx context.Context) (ListResponse, error) {
	resp, err := s.base.Client.V4().GetFortresses(
		s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ListResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetFortressesResponse)
	if err != nil {
		return ListResponse{ResponseMeta: meta}, err
	}

	return ListResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type Handle struct {
	client service.Client
	id     int
}

// Fortress returns a handle for a specific fortress with the given ID.
//
// Example:
//
//	fortress := client.Fortresses.Fortress(1)
//	_ = fortress
func (s *Service) Fortress(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

type Data = v4Client.FortressData

type InfoResponse struct {
	Data         Data
	ResponseMeta common.ResponseMeta
}

// Info retrieves detailed information about the fortress.
//
// It returns both the normalized Data and the raw API response.
//
// Example:
//
//	info, err := client.Fortresses.Fortress(1).Info(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Fortress: %s (ID: %d)\n", info.Data.Name, info.Data.Id)
func (h *Handle) Info(ctx context.Context) (InfoResponse, error) {
	resp, err := h.client.V4().GetFortress(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return InfoResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetFortressResponse)
	if err != nil {
		return InfoResponse{ResponseMeta: meta}, err
	}

	return InfoResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type SubmitFlagData struct {
	Message string
	Status  int
}

type SubmitFlagResponse struct {
	Data         SubmitFlagData
	ResponseMeta common.ResponseMeta
}

// SubmitFlag submits a flag for the fortress and returns the server's response message.
//
// Example:
//
//	result, err := client.Fortresses.Fortress(1).SubmitFlag(ctx, "HTB{3x4mp13_f14g}")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Flag submission: %s (Status: %d)\n", result.Data.Message, result.Data.Status)
func (h *Handle) SubmitFlag(ctx context.Context, flag string) (SubmitFlagResponse, error) {
	resp, err := h.client.V4().PostFortressFlag(
		h.client.Limiter().Wrap(ctx),
		h.id,
		v4Client.PostFortressFlagJSONRequestBody{
			Flag: flag,
		},
	)
	if err != nil {
		return SubmitFlagResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostFortressFlagResponse)
	if err != nil {
		return SubmitFlagResponse{ResponseMeta: meta}, err
	}

	return SubmitFlagResponse{
		Data: SubmitFlagData{
			Message: parsed.JSON200.Message,
			Status:  parsed.JSON200.Status,
		},
		ResponseMeta: meta,
	}, nil
}

type FlagData = common.FlagData

// Flags retrieves all available flags for the fortress.
//
// Example:
//
//	flags, err := client.Fortresses.Fortress(1).Flags(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Flags available: %d\n", len(flags.Flags))
func (h *Handle) Flags(ctx context.Context) (FlagData, error) {
	resp, err := h.client.V4().GetFortressFlags(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return FlagData{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetFortressFlagsResponse)
	if err != nil {
		return FlagData{ResponseMeta: meta}, err
	}

	return FlagData{
		Flags:        parsed.JSON200.Data,
		Status:       parsed.JSON200.Status,
		ResponseMeta: meta,
	}, nil
}

type ResetFlagData struct {
	Message string
	Status  bool
}

type ResetResponse struct {
	Data         ResetFlagData
	ResponseMeta common.ResponseMeta
}

// Reset sends a reset request for the associated fortress VM.
//
// Example:
//
//	result, err := client.Fortresses.Fortress(1).Reset(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Reset result: %s (Status: %t)\n", result.Data.Message, result.Data.Status)
func (h *Handle) Reset(ctx context.Context) (ResetResponse, error) {
	resp, err := h.client.V4().PostFortressReset(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return ResetResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostFortressResetResponse)
	if err != nil {
		return ResetResponse{ResponseMeta: meta}, err
	}

	return ResetResponse{
		Data: ResetFlagData{
			Message: parsed.JSON200.Message,
			Status:  parsed.JSON200.Status,
		},
		ResponseMeta: meta,
	}, nil
}
