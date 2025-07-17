package vms

import (
	"context"

	v4client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
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

// VM returns a handle for a specific virtual machine with the given ID.
// This handle can be used to perform operations on the VM such as
// spawning, resetting, extending, or terminating the instance.
// The ID is typically obtained from machine listings or other API responses.
func (s *Service) VM(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

// Reset performs a hard reset of the virtual machine.
// This operation forcefully restarts the VM instance.
//
// Example:
//
//	result, err := client.VMs.VM(12345).Reset(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Reset result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
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

// Spawn starts a new instance of the virtual machine.
// This creates and boots a VM instance for the specified machine.
//
// Example:
//
//	result, err := client.VMs.VM(12345).Spawn(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Spawn result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
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

// Extend extends the runtime of the virtual machine instance.
// This operation prolongs the active session time for the VM.
//
// Example:
//
//	result, err := client.VMs.VM(12345).Extend(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Extend result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
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

// Terminate stops and destroys the virtual machine instance.
// This operation permanently shuts down the VM and releases resources.
//
// Example:
//
//	result, err := client.VMs.VM(12345).Terminate(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Terminate result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
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

// VoteReset initiates a vote to reset the virtual machine.
// This is typically used in shared environments where multiple users
// can vote to reset a machine instance.
//
// Example:
//
//	result, err := client.VMs.VM(12345).VoteReset(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Vote reset result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
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

// VoteResetAccept accepts a pending reset vote for the virtual machine.
// This operation confirms participation in a reset vote that was
// previously initiated by VoteReset.
//
// Example:
//
//	result, err := client.VMs.VM(12345).VoteResetAccept(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Vote accept result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
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
