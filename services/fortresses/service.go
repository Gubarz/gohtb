package fortresses

import (
	"context"
	"fmt"
	"sort"

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
	resp, err := s.base.Client.V4().GetFortressesWithResponse(
		s.base.Client.Limiter().Wrap(ctx))

	raw := extract.Raw(resp)
	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) ListResponse {
			return ListResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}
	var list []Item
	if resp.JSON200.Data != nil {
		for _, fortress := range *resp.JSON200.Data {
			list = append(list, toItem(fortress))
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Id < list[j].Id
	})
	return ListResponse{
		Data: list,
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
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
	resp, err := h.client.V4().GetFortressWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) Info {
			return Info{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return Info{
		Data: fromFortressData(resp.JSON200.Data),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
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
	resp, err := h.client.V4().PostFortressFlagWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
		v4client.PostFortressFlagJSONRequestBody{
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
		Data: SubmitFlagData{
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
	resp, err := h.client.V4().GetFortressFlagsWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) FlagData {
			return FlagData{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return FlagData{
		Flags:  convert.Slice(*resp.JSON200.Data, common.FromAPIFlag),
		Status: deref.Bool(resp.JSON200.Status),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
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
	resp, err := h.client.V4().PostFortressResetWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return ResetResponse{}, fmt.Errorf("API error: %w", err)
	}
	return ResetResponse{
		Data: ResetFlagData{
			Message: deref.String(resp.JSON200.Message),
			Status:  deref.Bool(resp.JSON200.Status),
		},
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}
