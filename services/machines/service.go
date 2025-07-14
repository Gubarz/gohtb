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

// Active retrieves information about the currently active machine.
// This returns details about the machine that is currently available for solving.
//
// Example:
//
//	active, err := client.Machines.Active(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Active machine: %s (ID: %d)\n", active.Data.Name, active.Data.Id)
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

// Machine returns a handle for a specific machine with the given ID.
// This handle can be used to perform operations on the machine such as
// retrieving information, submitting flags, or managing VM instances.
func (s *Service) Machine(id int) *Handle {
	return &Handle{
		client:  s.base.Client,
		id:      id,
		product: s.product,
	}
}

// Info retrieves detailed information about the machine.
// This includes comprehensive machine details such as name, difficulty,
// operating system, release date, and other metadata.
//
// Example:
//
//	info, err := client.Machines.Machine(12345).Info(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Machine: %s (%s, %s)\n", info.Data.Name, info.Data.Os, info.Data.DifficultyText)
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

// Own submits a flag for the machine to claim ownership.
// This is used to submit user or root flags for machines to mark completion.
//
// Example:
//
//	result, err := client.Machines.Machine(12345).Own(ctx, "60b725f10c9c85c70d97880dfe8191b3")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Flag submission: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
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

// Reset performs a hard reset of the machine's virtual machine instance.
// This operation forcefully restarts the VM instance for this machine.
//
// Example:
//
//	result, err := client.Machines.Machine(12345).Reset(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Reset result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (h *Handle) Reset(ctx context.Context) (vms.Response, error) {
	return vms.NewService(h.client).VM(h.id).Reset(ctx)
}

// Extend extends the runtime of the machine's virtual machine instance.
// This operation prolongs the active session time for the VM.
//
// Example:
//
//	result, err := client.Machines.Machine(12345).Extend(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Extend result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (h *Handle) Extend(ctx context.Context) (vms.Response, error) {
	return vms.NewService(h.client).VM(h.id).Extend(ctx)
}

// Terminate stops and destroys the machine's virtual machine instance.
// This operation permanently shuts down the VM and releases resources.
//
// Example:
//
//	result, err := client.Machines.Machine(12345).Terminate(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Terminate result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (h *Handle) Terminate(ctx context.Context) (vms.Response, error) {
	return vms.NewService(h.client).VM(h.id).Terminate(ctx)
}

// Spawn starts a new instance of the machine's virtual machine.
// This creates and boots a VM instance for the specified machine.
//
// Example:
//
//	result, err := client.Machines.Machine(12345).Spawn(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Spawn result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (h *Handle) Spawn(ctx context.Context) (vms.Response, error) {
	return vms.NewService(h.client).VM(h.id).Spawn(ctx)
}

// ByOS filters machines by operating system using case-insensitive matching.
// Valid values include "Linux" and "Windows".
// Returns a new MachineDataItems slice containing only machines with the specified OS.
//
// Example:
//
//	linuxMachines := machines.Data.ByOS("Linux")
//	windowsMachines := machines.Data.ByOS("Windows")
func (m MachineDataItems) ByOS(os string) MachineDataItems {
	var d MachineDataItems
	for _, v := range m {
		if strings.EqualFold(v.Os, os) {
			d = append(d, v)
		}
	}
	return d
}

// First returns the first machine in the slice, or an empty MachineData if the slice is empty.
// This is commonly used after filtering to get the first match.
//
// Example:
//
//	firstLinux := machines.Data.ByOS("Linux").First()
func (m MachineDataItems) First() MachineData {
	if len(m) == 0 {
		return MachineData{}
	}
	return m[0]
}

// Last returns the last machine in the slice, or an empty MachineData if the slice is empty.
// This is commonly used to get the most recent or final machine in a list.
//
// Example:
//
//	lastMachine := machines.Data.Last()
func (m MachineDataItems) Last() MachineData {
	for i := len(m) - 1; i >= 0; i-- {
		return m[i]
	}
	return MachineData{}
}

// ByDifficulty filters machines by difficulty level using case-insensitive matching.
// Valid values are "Easy", "Medium", "Hard", and "Insane".
// Returns a new MachineDataItems slice containing only machines with the specified difficulty.
//
// Example:
//
//	hardMachines := machines.Data.ByDifficulty("Hard")
//	easyMachines := machines.Data.ByDifficulty("Easy")
func (m MachineDataItems) ByDifficulty(difficulty string) MachineDataItems {
	var out MachineDataItems
	for _, v := range m {
		if strings.EqualFold(v.DifficultyText, difficulty) {
			out = append(out, v)
		}
	}
	return out
}
