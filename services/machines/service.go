package machines

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	v5Client "github.com/gubarz/gohtb/httpclient/v5"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
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
type DifficultyChart = v4Client.DifficultyChart1

type MachinesData struct {
	v5Client.MachinesItem
}

type MachinesDataItems []MachinesData

type MachinesResponse struct {
	Data         MachinesDataItems
	Pagination   PagingMeta
	ResponseMeta common.ResponseMeta
}

type Service struct {
	base    service.Base
	product string
}

// NewService creates a new machines service bound to a shared client.
//
// Example:
//
//	machineService := machines.NewService(client, "machine")
//	_ = machineService
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
//
// Example:
//
//	machine := client.Machines.Machine(12345)
//	_ = machine
func (s *Service) Machine(id int) *Handle {
	return &Handle{
		client:  s.base.Client,
		id:      id,
		product: s.product,
	}
}

// MachineName returns a handle for a specific machine with the given slug/name.
//
// Example:
//
//	machine := client.Machines.MachineName("lame")
//	_ = machine
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
	IsAssumedBreach  bool
	FeedbackForChart DifficultyChart
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
	wrapped.FeedbackForChart = feedbackForChart(wrapped.MachineProfileInfo.FeedbackForChart)

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

type RecommendedMachinesData = v4Client.MachineRecommendedResponse

type RecommendedMachinesResponse struct {
	Data         RecommendedMachinesData
	ResponseMeta common.ResponseMeta
}

// Recommended retrieves currently recommended active machines.
//
// Example:
//
//	recommended, err := client.Machines.Recommended(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Recommended card 1: %s\n", recommended.Data.Card1.Name)
func (s *Service) Recommended(ctx context.Context) (RecommendedMachinesResponse, error) {
	resp, err := s.base.Client.V4().GetMachineRecommended(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return RecommendedMachinesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineRecommendedResponse)
	if err != nil {
		return RecommendedMachinesResponse{ResponseMeta: meta}, err
	}

	return RecommendedMachinesResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type WalkthroughRandomData = v4Client.MachineWalkthroughRandomResponse

type WalkthroughRandomResponse struct {
	Data         WalkthroughRandomData
	ResponseMeta common.ResponseMeta
}

// WalkthroughRandom returns a random machine walkthrough author.
//
// Example:
//
//	random, err := client.Machines.WalkthroughRandom(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Random author: %s\n", random.Data.Message.Name)
func (s *Service) WalkthroughRandom(ctx context.Context) (WalkthroughRandomResponse, error) {
	resp, err := s.base.Client.V4().GetMachineWalkthroughRandom(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return WalkthroughRandomResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineWalkthroughRandomResponse)
	if err != nil {
		return WalkthroughRandomResponse{ResponseMeta: meta}, err
	}

	return WalkthroughRandomResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type WalkthroughLanguage = v4Client.MachineWalkthroughsLanguageListItem

type WalkthroughLanguagesResponse struct {
	Data         []WalkthroughLanguage
	ResponseMeta common.ResponseMeta
}

// WalkthroughLanguages retrieves the available machine walkthrough languages.
//
// Example:
//
//	languages, err := client.Machines.WalkthroughLanguages(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Walkthrough languages: %d\n", len(languages.Data))
func (s *Service) WalkthroughLanguages(ctx context.Context) (WalkthroughLanguagesResponse, error) {
	resp, err := s.base.Client.V4().GetMachineWalkthroughsLanguageList(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return WalkthroughLanguagesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineWalkthroughsLanguageListResponse)
	if err != nil {
		return WalkthroughLanguagesResponse{ResponseMeta: meta}, err
	}

	return WalkthroughLanguagesResponse{
		Data:         parsed.JSON200.Info,
		ResponseMeta: meta,
	}, nil
}

type ActivityItem = v4Client.MachineAcitivtyItem

type ActivityResponse struct {
	Data         []ActivityItem
	ResponseMeta common.ResponseMeta
}

// Activity retrieves recent activity for the machine.
//
// Example:
//
//	activity, err := client.Machines.Machine(12345).Activity(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Activity items: %d\n", len(activity.Data))
func (h *Handle) Activity(ctx context.Context) (ActivityResponse, error) {
	resp, err := h.client.V4().GetMachineActivity(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return ActivityResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineActivityResponse)
	if err != nil {
		return ActivityResponse{ResponseMeta: meta}, err
	}

	return ActivityResponse{
		Data:         parsed.JSON200.Info.Activity,
		ResponseMeta: meta,
	}, nil
}

type ChangelogItem = v4Client.MachineChangeLogItem

type ChangelogResponse struct {
	Data         []ChangelogItem
	ResponseMeta common.ResponseMeta
}

// Changelog retrieves changelog entries for the machine.
//
// Example:
//
//	changelog, err := client.Machines.Machine(12345).Changelog(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Changelog entries: %d\n", len(changelog.Data))
func (h *Handle) Changelog(ctx context.Context) (ChangelogResponse, error) {
	resp, err := h.client.V4().GetMachineChangelog(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return ChangelogResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineChangelogResponse)
	if err != nil {
		return ChangelogResponse{ResponseMeta: meta}, err
	}

	return ChangelogResponse{
		Data:         parsed.JSON200.Info,
		ResponseMeta: meta,
	}, nil
}

type GraphMatrixData = v4Client.MachineGraphMatrixInfo

type GraphMatrixResponse struct {
	Data         GraphMatrixData
	ResponseMeta common.ResponseMeta
}

// GraphMatrix retrieves machine matrix graph data.
//
// Example:
//
//	matrix, err := client.Machines.Machine(12345).GraphMatrix(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("User enum score: %.2f\n", matrix.Data.User.Enum)
func (h *Handle) GraphMatrix(ctx context.Context) (GraphMatrixResponse, error) {
	resp, err := h.client.V4().GetMachineGraphMatrix(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return GraphMatrixResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineGraphMatrixResponse)
	if err != nil {
		return GraphMatrixResponse{ResponseMeta: meta}, err
	}

	return GraphMatrixResponse{
		Data:         parsed.JSON200.Info,
		ResponseMeta: meta,
	}, nil
}

type MachineTagItem = v4Client.MachineTagItems

type TagsResponse struct {
	Data         []MachineTagItem
	ResponseMeta common.ResponseMeta
}

// Tags retrieves tags assigned to the machine.
//
// Example:
//
//	tags, err := client.Machines.Machine(12345).Tags(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Machine tags: %d\n", len(tags.Data))
func (h *Handle) Tags(ctx context.Context) (TagsResponse, error) {
	resp, err := h.client.V4().GetMachineTags(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return TagsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineTagsResponse)
	if err != nil {
		return TagsResponse{ResponseMeta: meta}, err
	}

	return TagsResponse{
		Data:         parsed.JSON200.Info,
		ResponseMeta: meta,
	}, nil
}

type WalkthroughData = v4Client.MachineWalkthroughMessage

type WalkthroughsResponse struct {
	Data         WalkthroughData
	ResponseMeta common.ResponseMeta
}

// Walkthroughs retrieves walkthrough metadata for the machine.
//
// Example:
//
//	walkthroughs, err := client.Machines.Machine(12345).Walkthroughs(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Community writeups: %d\n", len(walkthroughs.Data.Writeups))
func (h *Handle) Walkthroughs(ctx context.Context) (WalkthroughsResponse, error) {
	resp, err := h.client.V4().GetMachineWalkthroughs(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return WalkthroughsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineWalkthroughsResponse)
	if err != nil {
		return WalkthroughsResponse{ResponseMeta: meta}, err
	}

	return WalkthroughsResponse{
		Data:         parsed.JSON200.Message,
		ResponseMeta: meta,
	}, nil
}

type WriteupResponse struct {
	Data         []byte
	ResponseMeta common.ResponseMeta
}

// Writeup downloads the machine writeup payload.
//
// Example:
//
//	writeup, err := client.Machines.Machine(12345).Writeup(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Writeup bytes: %d\n", len(writeup.Data))
func (h *Handle) Writeup(ctx context.Context) (WriteupResponse, error) {
	resp, err := h.client.V4().GetMachineWriteup(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) WriteupResponse {
			return WriteupResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return WriteupResponse{
		Data: raw,
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			CFRay:      resp.Header.Get("CF-Ray"),
		},
	}, nil
}

type AdventureData = v4Client.MachinesAdventureResponse

type AdventureResponse struct {
	Data         AdventureData
	ResponseMeta common.ResponseMeta
}

// Adventure retrieves machine adventure tasks.
//
// Example:
//
//	adventure, err := client.Machines.Machine(12345).Adventure(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Adventure tasks: %d\n", len(adventure.Data.Data))
func (h *Handle) Adventure(ctx context.Context) (AdventureResponse, error) {
	resp, err := h.client.V4().GetMachineAdventure(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return AdventureResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineAdventureResponse)
	if err != nil {
		return AdventureResponse{ResponseMeta: meta}, err
	}

	return AdventureResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type TasksData = v4Client.MachineTasksResponse

type TasksResponse struct {
	Data         TasksData
	ResponseMeta common.ResponseMeta
}

// Tasks retrieves machine task metadata.
//
// Example:
//
//	tasks, err := client.Machines.Machine(12345).Tasks(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Task entries: %d\n", len(tasks.Data.Data))
func (h *Handle) Tasks(ctx context.Context) (TasksResponse, error) {
	resp, err := h.client.V4().GetMachineTasks(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return TasksResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineTasksResponse)
	if err != nil {
		return TasksResponse{ResponseMeta: meta}, err
	}

	return TasksResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
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

type TagCategory = v4Client.TagCategory

type Tag = v4Client.Tag

type TagCatalog struct {
	Categories       []TagCategory
	Tags             []Tag
	TagsById         map[int]Tag
	CategoriesById   map[int]TagCategory
	TagsByCat        map[int][]Tag
	TagsByName       map[string]Tag
	CategoriesByName map[string]TagCategory
}

func feedbackForChart(u v4Client.DifficultyChart) DifficultyChart {
	n, err := u.AsDifficultyChart1()
	if err != nil {
		return DifficultyChart{}
	}
	return n
}
