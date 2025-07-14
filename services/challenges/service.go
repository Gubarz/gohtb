package challenges

import (
	"context"
	"strconv"

	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/deref"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	v4client "github.com/gubarz/gohtb/internal/httpclient/v4"
	"github.com/gubarz/gohtb/internal/service"
)

func NewService(client service.Client, product string) *Service {
	return &Service{
		base:    service.NewBase(client),
		product: product,
	}
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
	slug := strconv.Itoa(h.id)
	resp, err := h.client.V4().GetChallengeInfoWithResponse(
		h.client.Limiter().Wrap(ctx),
		slug,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil || resp.JSON200.Challenge == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) InfoResponse {
			return InfoResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}
	return InfoResponse{
		Data: fromAPIChallengeInfo(resp.JSON200.Challenge),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
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
	resp, err := h.client.V4().PostTodoUpdateWithResponse(
		h.client.Limiter().Wrap(ctx),
		v4client.PostTodoUpdateParamsProduct(h.product),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil || resp.JSON200.Info == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) common.TodoUpdateResponse {
			return common.TodoUpdateResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}
	return common.TodoUpdateResponse{
		Data: convert.SlicePointer(resp.JSON200.Info, common.FromAPIInfoArray),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
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
	resp, err := h.client.V4().PostChallengeStartWithFormdataBodyWithResponse(
		h.client.Limiter().Wrap(ctx),
		v4client.PostChallengeStartFormdataRequestBody{
			ChallengeId: h.id,
		},
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) common.MessageResponse {
			return common.MessageResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}
	return common.MessageResponse{
		Data: common.Message{
			Message: deref.String(resp.JSON200.Message),
		},
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
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
	resp, err := h.client.V4().PostChallengeStopWithFormdataBodyWithResponse(
		h.client.Limiter().Wrap(ctx),
		v4client.PostChallengeStopFormdataRequestBody{
			ChallengeId: h.id,
		},
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil || resp.JSON200.Message == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) common.MessageResponse {
			return common.MessageResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}
	return common.MessageResponse{
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
func (h *Handle) Own(ctx context.Context, flag string) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostChallengeOwnWithFormdataBodyWithResponse(
		h.client.Limiter().Wrap(ctx),
		v4client.ChallengeOwnRequest{
			ChallengeId: h.id,
			Flag:        flag,
		},
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil || resp.JSON200.Message == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) common.MessageResponse {
			return common.MessageResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}
	return common.MessageResponse{
		Data: common.Message{
			Message: deref.String(resp.JSON200.Message),
		},
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
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
	resp, err := h.client.V4().GetChallengeActivityWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil || resp.JSON200.Info == nil || resp.JSON200.Info.Activity == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) ActivityResponse {
			return ActivityResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}
	return ActivityResponse{
		Data: convert.SlicePointer(resp.JSON200.Info.Activity, fromAPIChallengeActivity),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
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
	resp, err := h.client.V4().GetChallengeDownloadWithResponse(
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
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}
