package fortresses

import (
	"context"
	"sort"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

// List retrieves basic information about the fortresses.
//
// It returns both the normalized Data and the raw API response.
//
// Example:
//
//	list, err := htb.Fortresses.List(ctx)
//	if err != nil {
//		if apiErr, ok := gohtb.AsAPIError(err); ok {
//		fmt.Println("API error:", apiErr.StatusCode, apiErr.Message)
//		} else {
//			fmt.Println("Unexpected error:", err)
//		}
//		return
//	}
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

	var list []Item
	if parsed.JSON200.Data != nil {
		for _, fortress := range parsed.JSON200.Data {
			list = append(list, fortress)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Id < list[j].Id
	})
	return ListResponse{
		Data:         list,
		ResponseMeta: meta,
	}, nil
}

func (s *Service) Fortress(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

// Info retrieves detailed information about the fortress.
//
// It returns both the normalized Data and the raw API response.
//
// Example:
//
//	info, err := htb.Fortresses.Fortress(1).Info(ctx)
//	if err != nil {
//		if apiErr, ok := gohtb.AsAPIError(err); ok {
//		fmt.Println("API error:", apiErr.StatusCode, apiErr.Message)
//		} else {
//			fmt.Println("Unexpected error:", err)
//		}
//		return
//	}
func (h *Handle) Info(ctx context.Context) (Info, error) {
	resp, err := h.client.V4().GetFortress(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return Info{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetFortressResponse)
	if err != nil {
		return Info{ResponseMeta: meta}, err
	}

	return Info{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

// SubmitFlag submits a flag for the fortress and returns the server's response message.
//
// Example:
//
//	msg, err := htb.Fortresses.Fortress(1).SubmitFlag(ctx, "HTB{3x4mp13_f14g}")
//	if err != nil {
//		if apiErr, ok := gohtb.AsAPIError(err); ok {
//		fmt.Println("API error:", apiErr.StatusCode, apiErr.Message)
//		} else {
//			fmt.Println("Unexpected error:", err)
//		}
//		return
//	}
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

// Flags retrieves all available flags for the fortress.
//
// Example:
//
//	flags, err := htb.Fortresses.Fortress(1).Flags(ctx)
//	if err != nil {
//		if apiErr, ok := gohtb.AsAPIError(err); ok {
//		fmt.Println("API error:", apiErr.StatusCode, apiErr.Message)
//		} else {
//			fmt.Println("Unexpected error:", err)
//		}
//		return
//	}
//	for _, flag := range flags {
//		fmt.Println(flag.Title)
//	}
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

// Reset sends a reset request for the associated fortress VM.
//
// Example:
//
//	msg, err := htb.Fortresses.Fortress(1).Reset(ctx)
//	if err != nil {
//		if apiErr, ok := gohtb.AsAPIError(err); ok {
//		fmt.Println("API error:", apiErr.StatusCode, apiErr.Message)
//		} else {
//			fmt.Println("Unexpected error:", err)
//		}
//		return
//	}
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
