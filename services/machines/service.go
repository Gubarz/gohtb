package machines

import (
	"context"
	"strconv"
	"strings"

	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	v5Client "github.com/gubarz/gohtb/internal/httpclient/v5"
	"github.com/gubarz/gohtb/internal/service"
	"github.com/gubarz/gohtb/services/vms"
)

func NewService(client service.Client, product string) *Service {
	return &Service{
		base:    service.NewBase(client),
		product: product,
	}
}

func (s *Service) Active(ctx context.Context) (ActiveResponse, error) {
	resp, err := s.base.Client.V4().GetMachineActiveWithResponse(s.base.Client.Limiter().Wrap(ctx))

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) ActiveResponse {
			return ActiveResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return ActiveResponse{
		Data: fromAPIActiveMachineInfo(resp.JSON200.Info),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (s *Service) Machine(id int) *Handle {
	return &Handle{
		client:  s.base.Client,
		id:      id,
		product: s.product,
	}
}

func (h *Handle) Info(ctx context.Context) (InfoResponse, error) {
	slug := strconv.Itoa(h.id)
	resp, err := h.client.V4().GetMachineProfileWithResponse(h.client.Limiter().Wrap(ctx), slug)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) InfoResponse {
			return InfoResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return InfoResponse{
		Data: fromAPIMachineProfileInfo(resp.JSON200.Info),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Own(ctx context.Context, flag string) (OwnResponse, error) {
	resp, err := h.client.V5().PostMachineOwnWithFormdataBodyWithResponse(h.client.Limiter().Wrap(ctx),
		v5Client.PostMachineOwnJSONRequestBody{
			Id:   h.id,
			Flag: flag,
		})

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) OwnResponse {
			return OwnResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return OwnResponse{
		Data: fromAPIMachineOwnResponse(*resp.JSON200),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

func (h *Handle) Reset(ctx context.Context) (vms.Response, error) {
	return vms.NewService(h.client).VM(h.id).Reset(ctx)
}

func (h *Handle) Extend(ctx context.Context) (vms.Response, error) {
	return vms.NewService(h.client).VM(h.id).Extend(ctx)
}

func (h *Handle) Terminate(ctx context.Context) (vms.Response, error) {
	return vms.NewService(h.client).VM(h.id).Terminate(ctx)
}

func (h *Handle) Spawn(ctx context.Context) (vms.Response, error) {
	return vms.NewService(h.client).VM(h.id).Spawn(ctx)
}

func (m MachineDataItems) ByOS(os string) MachineDataItems {
	var d MachineDataItems
	for _, v := range m {
		if strings.EqualFold(v.Os, os) {
			d = append(d, v)
		}
	}
	return d
}

func (m MachineDataItems) First() MachineData {
	if len(m) == 0 {
		return MachineData{}
	}
	return m[0]
}

func (m MachineDataItems) Last() MachineData {
	for i := len(m) - 1; i >= 0; i-- {
		return m[i]
	}
	return MachineData{}
}

func (m MachineDataItems) ByDifficulty(difficulty string) MachineDataItems {
	var out MachineDataItems
	for _, v := range m {
		if strings.EqualFold(v.DifficultyText, difficulty) {
			out = append(out, v)
		}
	}
	return out
}
