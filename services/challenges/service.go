package challenges

import (
	"context"
	"strconv"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	"github.com/gubarz/gohtb/internal/service"
	"github.com/gubarz/gohtb/services/containers"
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

// NewService creates a new challenges service bound to a shared client.
//
// Example:
//
//	challengeService := challenges.NewService(client, "challenge")
//	_ = challengeService
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

// Categories retrieves the challenge categories list.
//
// Example:
//
//	categories, err := client.Challenges.Categories(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Categories: %d\n", len(categories.Data))
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

type RecommendedData = v4Client.ChallengeRecommendedResponse
type RecommendedResponse struct {
	Data         RecommendedData
	ResponseMeta common.ResponseMeta
}

// Recommended retrieves currently recommended active challenges.
//
// Example:
//
//	recommended, err := client.Challenges.Recommended(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Recommended card 1: %s\n", recommended.Data.Card1.Name)
func (s *Service) Recommended(ctx context.Context) (RecommendedResponse, error) {
	resp, err := s.base.Client.V4().GetChallengeRecommended(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return RecommendedResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetChallengeRecommendedResponse)
	if err != nil {
		return RecommendedResponse{ResponseMeta: meta}, err
	}

	return RecommendedResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type RecommendedRetiredData = v4Client.ChallengeRecommendedRetiredResponse
type RecommendedRetiredResponse struct {
	Data         RecommendedRetiredData
	ResponseMeta common.ResponseMeta
}

// RecommendedRetired retrieves currently recommended retired challenges.
//
// Example:
//
//	retired, err := client.Challenges.RecommendedRetired(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Recommended retired card 1: %s\n", retired.Data.Card1.Name)
func (s *Service) RecommendedRetired(ctx context.Context) (RecommendedRetiredResponse, error) {
	resp, err := s.base.Client.V4().GetChallengeRecommendedRetired(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return RecommendedRetiredResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetChallengeRecommendedRetiredResponse)
	if err != nil {
		return RecommendedRetiredResponse{ResponseMeta: meta}, err
	}

	return RecommendedRetiredResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type SuggestedData = v4Client.ChallengeSuggestedResponse
type SuggestedResponse struct {
	Data         SuggestedData
	ResponseMeta common.ResponseMeta
}

// Suggested retrieves a suggested challenge for the current user.
//
// Example:
//
//	suggested, err := client.Challenges.Suggested(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Suggested challenge: %s\n", suggested.Data.Data.Name)
func (s *Service) Suggested(ctx context.Context) (SuggestedResponse, error) {
	resp, err := s.base.Client.V4().GetChallengeSuggested(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return SuggestedResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetChallengeSuggestedResponse)
	if err != nil {
		return SuggestedResponse{ResponseMeta: meta}, err
	}

	return SuggestedResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

// Challenge returns a handle for a specific challenge with the given ID.
// This handle can be used to perform operations on the challenge such as
// retrieving information, starting/stopping instances, or submitting flags.
//
// Example:
//
//	challenge := client.Challenges.Challenge(12345)
//	_ = challenge
func (s *Service) Challenge(id int) *Handle {
	return &Handle{
		client:  s.base.Client,
		id:      id,
		product: s.product,
	}
}

// ChallengeName returns a handle for a specific challenge with the given slug/name.
//
// Example:
//
//	challenge := client.Challenges.ChallengeName("example-challenge")
//	info, err := challenge.Info(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Challenge: %s\n", info.Data.Name)
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
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Challenges found: %d\n", len(challenges.Data))
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
	return containers.NewService(h.client).Container(h.id).Start(ctx)
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
	return containers.NewService(h.client).Container(h.id).Stop(ctx)
}

// Own submits a flag for the challenge to claim ownership.
// This is used to submit the solution flag to complete the challenge.
// The difficulty parameter is a user rating (1-100) of how difficult the challenge was.
// If difficulty is 0 or negative, it defaults to 10.
//
// Example:
//
//	result, err := client.Challenges.Challenge(12345).Own(ctx, "HTB{example_flag_here}", 10)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Flag submission: %s\n", result.Data.Message)
func (h *Handle) Own(ctx context.Context, flag string, difficulty int) (common.MessageResponse, error) {
	if difficulty <= 0 {
		difficulty = 10
	}
	resp, err := h.client.V4().PostChallengeOwnWithFormdataBody(
		h.client.Limiter().Wrap(ctx),
		v4Client.ChallengeOwnRequest{
			ChallengeId: h.id,
			Difficulty:  difficulty,
			Flag:        flag,
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

type ChangelogData = v4Client.ChallengeChangelogChallengeIdResponse
type ChangelogResponse struct {
	Data         ChangelogData
	ResponseMeta common.ResponseMeta
}

// Changelog retrieves changelog entries for the challenge.
//
// Example:
//
//	changelog, err := client.Challenges.Challenge(12345).Changelog(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Changelog entries: %d\n", len(changelog.Data.Data))
func (h *Handle) Changelog(ctx context.Context) (ChangelogResponse, error) {
	resp, err := h.client.V4().GetChallengeChangelog(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return ChangelogResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetChallengeChangelogResponse)
	if err != nil {
		return ChangelogResponse{ResponseMeta: meta}, err
	}

	return ChangelogResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type ReviewsUserData = v4Client.ChallengeReviewsUserChallengeIdResponse
type ReviewsUserResponse struct {
	Data         ReviewsUserData
	ResponseMeta common.ResponseMeta
}

// ReviewsUser retrieves the current user's review details for the challenge.
//
// Example:
//
//	review, err := client.Challenges.Challenge(12345).ReviewsUser(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Review info: %s\n", review.Data.Info)
func (h *Handle) ReviewsUser(ctx context.Context) (ReviewsUserResponse, error) {
	resp, err := h.client.V4().GetChallengeReviewsUser(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return ReviewsUserResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetChallengeReviewsUserResponse)
	if err != nil {
		return ReviewsUserResponse{ResponseMeta: meta}, err
	}

	return ReviewsUserResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type WriteupData = v4Client.WriteupData
type WriteupResponse struct {
	Data         WriteupData
	ResponseMeta common.ResponseMeta
}

// Writeup retrieves writeup metadata for the challenge.
//
// Example:
//
//	writeup, err := client.Challenges.Challenge(12345).Writeup(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Official writeup available: %t\n", writeup.Data.Official.Id != 0)
func (h *Handle) Writeup(ctx context.Context) (WriteupResponse, error) {
	resp, err := h.client.V4().GetChallengeWriteup(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return WriteupResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetChallengeWriteupResponse)
	if err != nil {
		return WriteupResponse{ResponseMeta: meta}, err
	}

	return WriteupResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type WriteupOfficialResponse struct {
	Data         []byte
	ResponseMeta common.ResponseMeta
}

// WriteupOfficial downloads the official writeup payload for the challenge.
//
// Example:
//
//	writeup, err := client.Challenges.Challenge(12345).WriteupOfficial(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Writeup bytes: %d\n", len(writeup.Data))
func (h *Handle) WriteupOfficial(ctx context.Context) (WriteupOfficialResponse, error) {
	resp, err := h.client.V4().GetChallengeWriteupOfficial(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) WriteupOfficialResponse {
			return WriteupOfficialResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return WriteupOfficialResponse{
		Data: raw,
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			CFRay:      resp.Header.Get("CF-Ray"),
		},
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
