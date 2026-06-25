package vms

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	v5Client "github.com/gubarz/gohtb/httpclient/v5"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

type Service struct {
	base service.Base
}

// NewService creates a new virtual machines service bound to a shared client.
//
// Example:
//
//	vmService := vms.NewService(client)
//	_ = vmService
func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

type Handle struct {
	client service.Client
	id     int
}

// VM returns a handle for a specific virtual machine with the given ID.
// This handle can be used to perform operations on the VM such as
// spawning, resetting, extending, or terminating the instance.
// The ID is typically obtained from machine listings or other API responses.
//
// Example:
//
//	vm := client.VMs.VM(12345)
//	_ = vm
func (s *Service) VM(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

type Response struct {
	Data         common.Message
	ResponseMeta common.ResponseMeta
}

type ActiveInfo = v5Client.VirtualMachineActiveInfo

type ActiveResponse struct {
	Data         ActiveInfo
	ResponseMeta common.ResponseMeta
}

// Active retrieves information about the currently active virtual machine.
//
// Example:
//
//	active, err := client.VMs.Active(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Active VM: %s (%s)\n", active.Data.Name, active.Data.Ip)
func (s *Service) Active(ctx context.Context) (ActiveResponse, error) {
	resp, err := s.base.Client.V5().GetVirtualMachineActive(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ActiveResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v5Client.ParseGetVirtualMachineActiveResponse)
	if err != nil {
		return ActiveResponse{ResponseMeta: meta}, err
	}

	rawInfo, err := parsed.JSON200.Info.MarshalJSON()
	if err != nil {
		return ActiveResponse{ResponseMeta: meta}, err
	}
	if string(rawInfo) == "null" {
		return ActiveResponse{ResponseMeta: meta}, nil
	}

	info, err := parsed.JSON200.Info.AsVirtualMachineActiveInfo()
	if err != nil {
		return ActiveResponse{ResponseMeta: meta}, err
	}

	return ActiveResponse{
		Data:         info,
		ResponseMeta: meta,
	}, nil
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
			Message: parsed.JSON200.Message,
			Success: parsed.JSON200.Success,
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
			Message: parsed.JSON200.Message,
			Success: parsed.JSON200.Success,
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
			Message: parsed.JSON200.Message,
			Success: parsed.JSON200.Success,
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
			Message: parsed.JSON200.Message,
			Success: parsed.JSON200.Success,
		},
		ResponseMeta: meta,
	}, nil
}
