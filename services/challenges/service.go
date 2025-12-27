package challenges

import (
	"context"
	"strconv"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	"github.com/gubarz/gohtb/internal/service"
)

type ChallengeQuery struct {
	client     service.Client
	status     v4Client.GetChallengesParamsStatus
	state      v4Client.State
	sortBy     v4Client.GetChallengesParamsSortBy
	sortType   v4Client.GetChallengesParamsSortType
	difficulty v4Client.Difficulty
	category   v4Client.Category
	keyword    v4Client.Keyword
	todo       v4Client.GetChallengesParamsTodo
	page       int
	perPage    int
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

type Handle struct {
	client  service.Client
	id      int
	name    string
	product string
}

type CategoriesListInfo = v4Client.CategoriesListInfo
type CategoriesListInfoResponse struct {
	Data         CategoriesListInfo
	ResponseMeta common.ResponseMeta
}

func (s *Service) Categories(ctx context.Context) (CategoriesListInfoResponse, error) {
	resp, err := s.base.Client.V4().GetChallengeCategoriesList(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return CategoriesListInfoResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetChallengeCategoriesListResponse)
	if err != nil {
		return CategoriesListInfoResponse{ResponseMeta: meta}, err
	}

	return CategoriesListInfoResponse{
		Data:         parsed.JSON200.Info,
		ResponseMeta: meta,
	}, nil
}

// Challenge returns a handle for a specific challenge with the given ID.
// This handle can be used to perform operations on the challenge such as
// retrieving information, starting/stopping instances, or submitting flags.
func (s *Service) Challenge(id int) *Handle {
	return &Handle{
		client:  s.base.Client,
		id:      id,
		product: s.product,
	}
}

func (s *Service) ChallengeName(name string) *Handle {
	return &Handle{
		client:  s.base.Client,
		name:    name,
		product: s.product,
	}
}

// List creates a new query for challenges.
// This returns a ChallengeQuery that can be chained with filtering and pagination methods.
// Use this to search and filter challenges based on various criteria.
//
// Example:
//
//	query := client.Challenges.List()
//	challenges, err := query.ByDifficulty("Hard").ByState("active").Results(ctx)
func (s *Service) List() *ChallengeQuery {
	return &ChallengeQuery{
		client:  s.base.Client,
		page:    1,
		perPage: 100,
	}
}

type InfoResponse struct {
	Data         Challenge
	ResponseMeta common.ResponseMeta
}

type DifficultyChart = v4Client.DifficultyChart1
type Points = v4Client.ChallengePoints0

type Challenge struct {
	v4Client.Challenge
	DifficultyChart v4Client.DifficultyChart1
	Points          v4Client.ChallengePoints0
}

func feedbackForChart(u v4Client.DifficultyChart) DifficultyChart {
	n, err := u.AsDifficultyChart1()
	if err != nil {
		return DifficultyChart{}
	}
	return n
}

func points(u v4Client.Challenge_Points) int {
	n, err := u.AsChallengePoints0()
	if err != nil {
		m, err := u.AsChallengePoints1()
		if err != nil {
			return 0
		}
		v, err := strconv.Atoi(m)
		if err != nil {
			return 0
		}
		return v
	}
	return n
}

func wrapChallengeInfo(x v4Client.Challenge) Challenge {
	return Challenge{Challenge: x}
}

// Info retrieves detailed information about the challenge.
// This includes comprehensive challenge details such as name, difficulty,
// category, description, and other metadata.
//
// Example:
//
//	info, err := client.Challenges.Challenge(12345).Info(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Challenge: %s (%s, %s)\n", info.Data.Name, info.Data.Difficulty, info.Data.Category)
func (h *Handle) Info(ctx context.Context) (InfoResponse, error) {
	var slug string
	if h.name != "" {
		slug = h.name
	} else {
		slug = strconv.Itoa(h.id)
	}
	resp, err := h.client.V4().GetChallengeInfo(
		h.client.Limiter().Wrap(ctx),
		slug,
	)
	if err != nil {
		return InfoResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetChallengeInfoResponse)
	if err != nil {
		return InfoResponse{ResponseMeta: meta}, err
	}

	wrapped := wrapChallengeInfo(parsed.JSON200.Challenge)
	wrapped.DifficultyChart = feedbackForChart(wrapped.Challenge.DifficultyChart)
	wrapped.Points = points(wrapped.Challenge.Points)

	return InfoResponse{
		Data:         wrapped,
		ResponseMeta: meta,
	}, nil
}

// ToDo toggles the challenge's todo status for the current user.
// This adds or removes the challenge from the user's todo list.
//
// Example:
//
//	result, err := client.Challenges.Challenge(12345).ToDo(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Todo updated: %+v\n", result.Data)
func (h *Handle) ToDo(ctx context.Context) (common.TodoUpdateResponse, error) {
	resp, err := h.client.V4().PostTodoUpdate(
		h.client.Limiter().Wrap(ctx),
		v4Client.PostTodoUpdateParamsProduct(h.product),
		h.id,
	)
	if err != nil {
		return common.TodoUpdateResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostTodoUpdateResponse)
	if err != nil {
		return common.TodoUpdateResponse{ResponseMeta: meta}, err
	}

	return common.TodoUpdateResponse{
		Data:         parsed.JSON200.Info,
		ResponseMeta: meta,
	}, nil
}

// Start initiates a new instance of the challenge.
// This creates and starts a challenge instance for solving.
//
// Example:
//
//	result, err := client.Challenges.Challenge(12345).Start(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Challenge started: %s\n", result.Data.Message)
func (h *Handle) Start(ctx context.Context) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostChallengeStartWithFormdataBody(
		h.client.Limiter().Wrap(ctx),
		v4Client.PostChallengeStartFormdataRequestBody{
			ChallengeId: h.id,
		},
	)
	if err != nil {
		return common.MessageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostChallengeStartResponse)
	if err != nil {
		return common.MessageResponse{ResponseMeta: meta}, err
	}

	return common.MessageResponse{
		Data: common.Message{
			Message: parsed.JSON200.Message,
		},
		ResponseMeta: meta,
	}, nil
}

// Stop terminates the running challenge instance.
// This stops and destroys the active challenge instance.
//
// Example:
//
//	result, err := client.Challenges.Challenge(12345).Stop(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Challenge stopped: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (h *Handle) Stop(ctx context.Context) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostChallengeStopWithFormdataBody(
		h.client.Limiter().Wrap(ctx),
		v4Client.PostChallengeStopFormdataRequestBody{
			ChallengeId: h.id,
		},
	)
	if err != nil {
		return common.MessageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostChallengeStopResponse)
	if err != nil {
		return common.MessageResponse{ResponseMeta: meta}, err
	}

	return common.MessageResponse{
		Data: common.Message{
			Message: parsed.JSON200.Message,
			Success: parsed.JSON200.Success,
		},
		ResponseMeta: meta,
	}, nil
}

// Own submits a flag for the challenge to claim ownership.
// This is used to submit the solution flag to complete the challenge.
//
// Example:
//
//	result, err := client.Challenges.Challenge(12345).Own(ctx, "HTB{example_flag_here}")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Flag submission: %s\n", result.Data.Message)
func (h *Handle) Own(ctx context.Context, flag string, difficulty int) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostChallengeOwnWithFormdataBody(
		h.client.Limiter().Wrap(ctx),
		v4Client.ChallengeOwnRequest{
			ChallengeId: h.id,
			Flag:        flag,
			Difficulty:  difficulty,
		},
	)
	if err != nil {
		return common.MessageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostChallengeOwnResponse)
	if err != nil {
		return common.MessageResponse{ResponseMeta: meta}, err
	}

	return common.MessageResponse{
		Data: common.Message{
			Message: parsed.JSON200.Message,
		},
		ResponseMeta: meta,
	}, nil
}

type ChallengeActivity = v4Client.ChallengeActivity

type ActivityResponse struct {
	Data         []ChallengeActivity
	ResponseMeta common.ResponseMeta
}

// Activity retrieves the recent activity history for the challenge.
// This includes recent solves, attempts, and other challenge-related activities.
//
// Example:
//
//	activity, err := client.Challenges.Challenge(12345).Activity(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, act := range activity.Data {
//		fmt.Printf("Activity: %s at %s\n", act.Type, act.Date)
//	}
func (h *Handle) Activity(ctx context.Context) (ActivityResponse, error) {
	resp, err := h.client.V4().GetChallengeActivity(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return ActivityResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetChallengeActivityResponse)
	if err != nil {
		return ActivityResponse{ResponseMeta: meta}, err
	}

	return ActivityResponse{
		Data:         parsed.JSON200.Info.Activity,
		ResponseMeta: meta,
	}, nil
}

type DownloadResponse struct {
	Data         []byte
	ResponseMeta common.ResponseMeta
}

// Download retrieves the challenge files for download.
// This returns the challenge's downloadable zip.
// Note: Not all challenges have downloadable files. Check Info() first to verify availability.
//
// Example:
//
//	challenge := client.Challenges.Challenge(12345)
//
//	// Get challenge info first to check if files are available
//	info, err := challenge.Info(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Check if the challenge has downloadable files
//	if info.Data.FileName == "" {
//		fmt.Println("This challenge has no downloadable files")
//		return
//	}
//
//	// Download the challenge files
//	download, err := challenge.Download(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Write to file
//	err = os.WriteFile(info.Data.FileName, download.Data, 0644)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Downloaded challenge files to: %s (%d bytes)\n", info.Data.FileName, len(download.Data))
func (h *Handle) Download(ctx context.Context) (DownloadResponse, error) {
	resp, err := h.client.V4().GetChallengeDownload(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) DownloadResponse {
			return DownloadResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}
	return DownloadResponse{
		Data: raw,
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			CFRay:      resp.Header.Get("CF-Ray"),
		},
	}, nil
}
