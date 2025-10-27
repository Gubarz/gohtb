package machines

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	v5Client "github.com/gubarz/gohtb/httpclient/v5"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
	"github.com/gubarz/gohtb/services/vms"
)

type PagingMeta struct {
	CurrentPage int
	PerPage     int
	Total       int
	TotalPages  int
	Count       int
}

type MachineData = v4Client.MachineData

type MachineDataItems []MachineData

type MachinePaginatedResponse struct {
	Data         MachineDataItems
	Pagination   PagingMeta
	ResponseMeta common.ResponseMeta
}

type Service struct {
	base    service.Base
	product string
}

func NewService(client service.Client, product string) *Service {
	return &Service{
		base:    service.NewBase(client),
		product: product,
	}
}

type ActiveMachineInfo = v4Client.ActiveMachineInfo

type ActiveResponse struct {
	Data         ActiveMachineInfo
	ResponseMeta common.ResponseMeta
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
	resp, err := s.base.Client.V4().GetMachineActive(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ActiveResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineActiveResponse)
	if err != nil {
		return ActiveResponse{ResponseMeta: meta}, err
	}

	return ActiveResponse{
		Data:         parsed.JSON200.Info,
		ResponseMeta: meta,
	}, nil
}

type Handle struct {
	client  service.Client
	id      int
	name    string
	product string
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

func (s *Service) MachineName(name string) *Handle {
	return &Handle{
		client:  s.base.Client,
		name:    name,
		product: s.product,
	}
}

type Credentials struct {
	Username string
	Password string
}

type MachineProfileInfo struct {
	v4Client.MachineProfileInfo
	Credentials
	IsAssumedBreach bool
}

type InfoResponse struct {
	Data         MachineProfileInfo
	ResponseMeta common.ResponseMeta
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
	var slug string
	if h.name != "" {
		slug = h.name
	} else {
		slug = strconv.Itoa(h.id)
	}
	resp, err := h.client.V4().GetMachineProfile(h.client.Limiter().Wrap(ctx), slug)

	if err != nil {
		return InfoResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineProfileResponse)
	if err != nil {
		return InfoResponse{ResponseMeta: meta}, err
	}

	wrapped := wrapMachineProfileInfo(parsed.JSON200.Info)
	wrapped.IsAssumedBreach, wrapped.Credentials = parseAssumedBreachStatus(wrapped.InfoStatus)

	return InfoResponse{
		Data:         wrapped,
		ResponseMeta: meta,
	}, nil
}

type MachineOwnResponse = v5Client.MachineOwnResponse

type OwnResponse struct {
	Data         MachineOwnResponse
	ResponseMeta common.ResponseMeta
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
	resp, err := h.client.V5().PostMachineOwnWithFormdataBody(h.client.Limiter().Wrap(ctx),
		v5Client.PostMachineOwnJSONRequestBody{
			Id:   h.id,
			Flag: flag,
		})

	if err != nil {
		return OwnResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v5Client.ParsePostMachineOwnResponse)
	if err != nil {
		return OwnResponse{ResponseMeta: meta}, err
	}

	return OwnResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
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
	if len(m) == 0 {
		return MachineData{}
	}
	return m[len(m)-1]
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

func parseAssumedBreachStatus(info string) (bool, Credentials) {
	if strings.Contains(strings.ToLower(info), "as is common in") {
		return true, extractCredentials(info)
	}
	return false, Credentials{}
}

func extractCredentials(s string) Credentials {
	// 1) Username: <user> Password: <pass>
	reUserPass := regexp.MustCompile(`(?i)\busername\b[:\s]*([^\s:;,/]+)\s*(?:[:;,\-]\s*)?\bpassword\b[:\s]*(\S+)`)
	if m := reUserPass.FindStringSubmatch(s); len(m) == 3 {
		return Credentials{Username: strings.TrimSpace(m[1]), Password: strings.TrimSpace(m[2])}
	}

	// 2) account ...: user / pass   (covers the MSSQL example and similar)
	reAccount := regexp.MustCompile(`(?i)account[^:]*:\s*([^\s/]+)\s*/\s*(\S+)`)
	if m := reAccount.FindStringSubmatch(s); len(m) == 3 {
		return Credentials{Username: strings.TrimSpace(m[1]), Password: strings.TrimSpace(m[2])}
	}

	// 3) fallback: "user X / Y" like "user: alice / p@ss"
	reLoose := regexp.MustCompile(`(?i)\b(?:user|username|account)\b[:\s]*([^\s/]+)\s*/\s*(\S+)`)
	if m := reLoose.FindStringSubmatch(s); len(m) == 3 {
		return Credentials{Username: strings.TrimSpace(m[1]), Password: strings.TrimSpace(m[2])}
	}

	return Credentials{}
}

func wrapMachineProfileInfo(x v4Client.MachineProfileInfo) MachineProfileInfo {
	return MachineProfileInfo{MachineProfileInfo: x}
}
