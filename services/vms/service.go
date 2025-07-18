package vms

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/deref"
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
	params := v4Client.PostVMResetJSONRequestBody{
		MachineId: h.id,
	}

	resp, err := h.client.V4().PostVMReset(h.client.Limiter().Wrap(ctx), params)
	if err != nil {
		return Response{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostVMResetResponse)
	if err != nil {
		return Response{ResponseMeta: meta}, err
	}

	return Response{
		Data: common.Message{
			Message: deref.String(parsed.JSON200.Message),
			Success: deref.Bool(parsed.JSON200.Success),
		},
		ResponseMeta: meta,
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
	params := v4Client.PostVMSpawnJSONRequestBody{
		MachineId: h.id,
	}
	resp, err := h.client.V4().PostVMSpawn(h.client.Limiter().Wrap(ctx), params)
	if err != nil {
		return Response{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostVMSpawnResponse)
	if err != nil {
		return Response{ResponseMeta: meta}, err
	}

	return Response{
		Data: common.Message{
			Message: deref.String(parsed.JSON200.Message),
			Success: deref.Bool(parsed.JSON200.Success),
		},
		ResponseMeta: meta,
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
	params := v4Client.PostVMExtendJSONRequestBody{
		MachineId: h.id,
	}
	resp, err := h.client.V4().PostVMExtend(h.client.Limiter().Wrap(ctx), params)
	if err != nil {
		return Response{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostVMExtendResponse)
	if err != nil {
		return Response{ResponseMeta: meta}, err
	}

	return Response{
		Data: common.Message{
			Message: deref.String(parsed.JSON200.Message),
			Success: deref.Bool(parsed.JSON200.Success),
		},
		ResponseMeta: meta,
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
	params := v4Client.PostVMTerminateJSONRequestBody{
		MachineId: h.id,
	}
	req, err := h.client.V4().PostVMTerminate(h.client.Limiter().Wrap(ctx), params)
	if err != nil {
		return Response{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(req, v4Client.ParsePostVMTerminateResponse)
	if err != nil {
		return Response{ResponseMeta: meta}, err
	}

	return Response{
		Data: common.Message{
			Message: deref.String(parsed.JSON200.Message),
			Success: deref.Bool(parsed.JSON200.Success),
		},
		ResponseMeta: meta,
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
	params := v4Client.PostVMResetVoteJSONRequestBody{
		MachineId: h.id,
	}
	resp, err := h.client.V4().PostVMResetVote(h.client.Limiter().Wrap(ctx), params)
	if err != nil {
		return Response{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostVMResetVoteResponse)
	if err != nil {
		return Response{ResponseMeta: meta}, err
	}

	return Response{
		Data: common.Message{
			Message: deref.String(parsed.JSON200.Message),
			Success: deref.Bool(parsed.JSON200.Success),
		},
		ResponseMeta: meta,
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
	params := v4Client.PostVMResetVoteAcceptJSONRequestBody{
		MachineId: h.id,
	}
	resp, err := h.client.V4().PostVMResetVoteAccept(h.client.Limiter().Wrap(ctx), params)
	if err != nil {
		return Response{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostVMResetVoteAcceptResponse)
	if err != nil {
		return Response{ResponseMeta: meta}, err
	}

	return Response{
		Data: common.Message{
			Message: deref.String(parsed.JSON200.Message),
			Success: deref.Bool(parsed.JSON200.Success),
		},
		ResponseMeta: meta,
	}, nil
}
