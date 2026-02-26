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

type MachineData struct {
	v4Client.MachineData
	IsTodo           bool
	FeedbackForChart v4Client.DifficultyChart1
}

type MachineDataItems []MachineData

type MachinePaginatedResponse struct {
	Data         MachineDataItems
	Pagination   PagingMeta
	ResponseMeta common.ResponseMeta
}

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

type RecommendedRetiredMachinesData = v4Client.MachineRecommendedRetiredResponse

type RecommendedRetiredMachinesResponse struct {
	Data         RecommendedRetiredMachinesData
	ResponseMeta common.ResponseMeta
}

// RecommendedRetired retrieves currently recommended retired machines.
//
// Example:
//
//	retired, err := client.Machines.RecommendedRetired(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Recommended retired card 1: %s\n", retired.Data.Card1.Name)
func (s *Service) RecommendedRetired(ctx context.Context) (RecommendedRetiredMachinesResponse, error) {
	resp, err := s.base.Client.V4().GetMachineRecommendedRetired(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return RecommendedRetiredMachinesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineRecommendedRetiredResponse)
	if err != nil {
		return RecommendedRetiredMachinesResponse{ResponseMeta: meta}, err
	}

	return RecommendedRetiredMachinesResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type UnreleasedMachinesData = v4Client.MachineUnreleasedResponse

type UnreleasedMachinesResponse struct {
	Data         UnreleasedMachinesData
	ResponseMeta common.ResponseMeta
}

// Unreleased retrieves unreleased machines.
//
// Example:
//
//	unreleased, err := client.Machines.Unreleased(ctx, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Unreleased machines: %d\n", len(unreleased.Data.Data))
func (s *Service) Unreleased(ctx context.Context, params *v4Client.GetMachineUnreleasedParams) (UnreleasedMachinesResponse, error) {
	resp, err := s.base.Client.V4().GetMachineUnreleased(s.base.Client.Limiter().Wrap(ctx), params)
	if err != nil {
		return UnreleasedMachinesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineUnreleasedResponse)
	if err != nil {
		return UnreleasedMachinesResponse{ResponseMeta: meta}, err
	}

	return UnreleasedMachinesResponse{
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

type WalkthroughOfficialFeedbackChoicesData = v4Client.MachineWalkthroughOfficialFeedbackChoicesResponse

type WalkthroughOfficialFeedbackChoicesResponse struct {
	Data         WalkthroughOfficialFeedbackChoicesData
	ResponseMeta common.ResponseMeta
}

// WalkthroughOfficialFeedbackChoices retrieves official walkthrough feedback options.
//
// Example:
//
//	choices, err := client.Machines.WalkthroughOfficialFeedbackChoices(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Feedback choices: %d\n", len(choices.Data.FeedbackChoices))
func (s *Service) WalkthroughOfficialFeedbackChoices(ctx context.Context) (WalkthroughOfficialFeedbackChoicesResponse, error) {
	resp, err := s.base.Client.V4().GetMachineWalkthroughOfficialFeedbackChoices(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return WalkthroughOfficialFeedbackChoicesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineWalkthroughOfficialFeedbackChoicesResponse)
	if err != nil {
		return WalkthroughOfficialFeedbackChoicesResponse{ResponseMeta: meta}, err
	}

	return WalkthroughOfficialFeedbackChoicesResponse{
		Data:         *parsed.JSON200,
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

type CreatorsData = v4Client.MachineCreatorResponse

type CreatorsResponse struct {
	Data         CreatorsData
	ResponseMeta common.ResponseMeta
}

// Creators retrieves creator and co-creator information for the machine.
//
// Example:
//
//	creators, err := client.Machines.Machine(12345).Creators(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Creator records: %d\n", len(creators.Data.Creator))
func (h *Handle) Creators(ctx context.Context) (CreatorsResponse, error) {
	resp, err := h.client.V4().GetMachineCreators(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return CreatorsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineCreatorsResponse)
	if err != nil {
		return CreatorsResponse{ResponseMeta: meta}, err
	}

	return CreatorsResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type GraphPeriod = v4Client.GetMachineGraphActivityParamsPeriod

const (
	GraphPeriod1W GraphPeriod = v4Client.GetMachineGraphActivityParamsPeriodN1W
	GraphPeriod1M GraphPeriod = v4Client.GetMachineGraphActivityParamsPeriodN1M
	GraphPeriod3M GraphPeriod = v4Client.GetMachineGraphActivityParamsPeriodN3M
	GraphPeriod6M GraphPeriod = v4Client.GetMachineGraphActivityParamsPeriodN6M
	GraphPeriod1Y GraphPeriod = v4Client.GetMachineGraphActivityParamsPeriodN1Y
)

type GraphActivityData = v4Client.MachineGraphActivityInfo

type GraphActivityResponse struct {
	Data         GraphActivityData
	ResponseMeta common.ResponseMeta
}

// GraphActivity retrieves machine ownership activity graph data for a period.
//
// Example:
//
//	graph, err := client.Machines.Machine(12345).GraphActivity(ctx, machines.GraphPeriod1M)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Graph points: %d\n", len(graph.Data.Periods))
func (h *Handle) GraphActivity(ctx context.Context, period GraphPeriod) (GraphActivityResponse, error) {
	if period == "" {
		period = GraphPeriod1M
	}

	resp, err := h.client.V4().GetMachineGraphActivity(
		h.client.Limiter().Wrap(ctx),
		h.id,
		v4Client.GetMachineGraphActivityParamsPeriod(period),
	)
	if err != nil {
		return GraphActivityResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineGraphActivityResponse)
	if err != nil {
		return GraphActivityResponse{ResponseMeta: meta}, err
	}

	return GraphActivityResponse{
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

type GraphOwnsDifficultyData = v4Client.MachineGraphOwnsDifficultyInfoItems

type GraphOwnsDifficultyResponse struct {
	Data         GraphOwnsDifficultyData
	ResponseMeta common.ResponseMeta
}

// GraphOwnsDifficulty retrieves ownership distribution by machine difficulty.
//
// Example:
//
//	difficulty, err := client.Machines.Machine(12345).GraphOwnsDifficulty(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Difficulty 1 user owns: %d\n", difficulty.Data.N1.User)
func (h *Handle) GraphOwnsDifficulty(ctx context.Context) (GraphOwnsDifficultyResponse, error) {
	resp, err := h.client.V4().GetMachineGraphOwnsDifficulty(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return GraphOwnsDifficultyResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineGraphOwnsDifficultyResponse)
	if err != nil {
		return GraphOwnsDifficultyResponse{ResponseMeta: meta}, err
	}

	return GraphOwnsDifficultyResponse{
		Data:         parsed.JSON200.Info,
		ResponseMeta: meta,
	}, nil
}

type OwnsTopItem = v4Client.MachineOwnsTopItem

type OwnsTopResponse struct {
	Data         []OwnsTopItem
	ResponseMeta common.ResponseMeta
}

// OwnsTop retrieves top ownership records for the machine.
//
// Example:
//
//	top, err := client.Machines.Machine(12345).OwnsTop(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Top owns entries: %d\n", len(top.Data))
func (h *Handle) OwnsTop(ctx context.Context) (OwnsTopResponse, error) {
	resp, err := h.client.V4().GetMachineOwnsTop(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return OwnsTopResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineOwnsTopResponse)
	if err != nil {
		return OwnsTopResponse{ResponseMeta: meta}, err
	}

	return OwnsTopResponse{
		Data:         parsed.JSON200.Info,
		ResponseMeta: meta,
	}, nil
}

type ReviewRequest struct {
	Headline string
	Review   string
	Stars    float32
}

type ReviewData = v4Client.MachineReviewResponse

type ReviewResponse struct {
	Data         ReviewData
	ResponseMeta common.ResponseMeta
}

// Review submits a machine review for the authenticated user.
//
// Example:
//
//	review, err := client.Machines.Machine(12345).Review(ctx, machines.ReviewRequest{
//		Headline: "Great box",
//		Review:   "Good progression and realistic chain.",
//		Stars:    4,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Review status entries: %d\n", len(review.Data.Message))
func (h *Handle) Review(ctx context.Context, req ReviewRequest) (ReviewResponse, error) {
	resp, err := h.client.V4().PostMachineReviewWithFormdataBody(
		h.client.Limiter().Wrap(ctx),
		v4Client.PostMachineReviewFormdataRequestBody{
			Id:       h.id,
			Headline: req.Headline,
			Review:   req.Review,
			Stars:    req.Stars,
		},
	)
	if err != nil {
		return ReviewResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostMachineReviewResponse)
	if err != nil {
		return ReviewResponse{ResponseMeta: meta}, err
	}

	return ReviewResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type ReviewsData = v4Client.MachineReviewsResponse

type ReviewsResponse struct {
	Data         ReviewsData
	ResponseMeta common.ResponseMeta
}

// Reviews retrieves machine reviews.
//
// Example:
//
//	reviews, err := client.Machines.Machine(12345).Reviews(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Reviews count: %.0f\n", reviews.Data.Count)
func (h *Handle) Reviews(ctx context.Context) (ReviewsResponse, error) {
	resp, err := h.client.V4().GetMachineReviews(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return ReviewsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineReviewsResponse)
	if err != nil {
		return ReviewsResponse{ResponseMeta: meta}, err
	}

	return ReviewsResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type ReviewsUserData = v4Client.MachineReviewsUserMachineIdResponse

type ReviewsUserResponse struct {
	Data         ReviewsUserData
	ResponseMeta common.ResponseMeta
}

// ReviewsUser retrieves the authenticated user's review status for the machine.
//
// Example:
//
//	review, err := client.Machines.Machine(12345).ReviewsUser(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Review status: %s\n", review.Data.Message)
func (h *Handle) ReviewsUser(ctx context.Context) (ReviewsUserResponse, error) {
	resp, err := h.client.V4().GetMachineReviewsUser(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return ReviewsUserResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetMachineReviewsUserResponse)
	if err != nil {
		return ReviewsUserResponse{ResponseMeta: meta}, err
	}

	return ReviewsUserResponse{
		Data:         *parsed.JSON200,
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
