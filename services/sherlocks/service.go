package sherlocks

import (
	"context"
	"strconv"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	"github.com/gubarz/gohtb/internal/service"
)

type SherlockQuery struct {
	client     service.Client
	status     v4Client.GetSherlocksParamsStatus
	state      v4Client.State
	sortBy     v4Client.GetSherlocksParamsSortBy
	sortType   v4Client.GetSherlocksParamsSortType
	difficulty v4Client.Difficulty
	category   v4Client.Category
	keyword    v4Client.Keyword
	page       int
	perPage    int
}

type Service struct {
	base service.Base
}

// NewService creates a new sherlocks service bound to a shared client.
//
// Example:
//
//	sherlockService := sherlocks.NewService(client)
//	_ = sherlockService
func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

type Handle struct {
	client service.Client
	id     int
}

type CategoriesListInfo = v4Client.CategoriesListInfo

type CategoriesListInfoResponse struct {
	Data         CategoriesListInfo
	ResponseMeta common.ResponseMeta
}

// Categories retrieves the sherlock category list.
//
// Example:
//
//	categories, err := client.Sherlocks.Categories(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sherlock categories: %d\n", len(categories.Data))
func (s *Service) Categories(ctx context.Context) (CategoriesListInfoResponse, error) {
	resp, err := s.base.Client.V4().GetSherlocksCategoriesList(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return CategoriesListInfoResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSherlocksCategoriesListResponse)
	if err != nil {
		return CategoriesListInfoResponse{ResponseMeta: meta}, err
	}

	return CategoriesListInfoResponse{
		Data:         parsed.JSON200.Info,
		ResponseMeta: meta,
	}, nil
}

// Sherlock returns a handle for a specific sherlock challenge by ID.
//
// Example:
//
//	sherlock := client.Sherlocks.Sherlock(123)
//	_ = sherlock
func (s *Service) Sherlock(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

// List creates a new query builder for sherlock listings.
//
// Example:
//
//	query := client.Sherlocks.List()
//	_ = query
func (s *Service) List() *SherlockQuery {
	return &SherlockQuery{
		client:  s.base.Client,
		page:    1,
		perPage: 100,
	}
}

type SherlockNamedItemData = v4Client.SherlockNamedItemData

type InfoResponse struct {
	Data         SherlockNamedItemData
	ResponseMeta common.ResponseMeta
}

// Info retrieves basic sherlock details.
//
// Example:
//
//	info, err := client.Sherlocks.Sherlock(123).Info(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sherlock: %s\n", info.Data.Name)
func (h *Handle) Info(ctx context.Context) (InfoResponse, error) {
	slug := strconv.Itoa(h.id)
	resp, err := h.client.V4().GetSherlock(h.client.Limiter().Wrap(ctx), slug)

	if err != nil {
		return InfoResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSherlockResponse)
	if err != nil {
		return InfoResponse{ResponseMeta: meta}, err
	}

	return InfoResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type SherlockPlay = v4Client.SherlockPlay

type PlayResponse struct {
	Data         SherlockPlay
	ResponseMeta common.ResponseMeta
}

// Play starts or resumes a sherlock play session.
//
// Example:
//
//	play, err := client.Sherlocks.Sherlock(123).Play(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Play status: %s\n", play.Data.Status)
func (h *Handle) Play(ctx context.Context) (PlayResponse, error) {
	resp, err := h.client.V4().GetSherlockPlay(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return PlayResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSherlockPlayResponse)
	if err != nil {
		return PlayResponse{ResponseMeta: meta}, err
	}

	return PlayResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type SherlockDownloadLink = v4Client.SherlockDownloadLink

type DownloadResponse struct {
	Data         SherlockDownloadLink
	ResponseMeta common.ResponseMeta
}

// DownloadLink retrieves the downloadable asset link for the sherlock.
//
// Example:
//
//	download, err := client.Sherlocks.Sherlock(123).DownloadLink(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Download URL: %s\n", download.Data.Link)
func (h *Handle) DownloadLink(ctx context.Context) (DownloadResponse, error) {
	resp, err := h.client.V4().GetSherlockDownloadlink(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return DownloadResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSherlockDownloadlinkResponse)
	if err != nil {
		return DownloadResponse{ResponseMeta: meta}, err
	}

	return DownloadResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type SherlockProgressData = v4Client.SherlockProgressData

type ProgressResponse struct {
	Data         SherlockProgressData
	ResponseMeta common.ResponseMeta
}

// Progress retrieves the authenticated user's progress for the sherlock.
//
// Example:
//
//	progress, err := client.Sherlocks.Sherlock(123).Progress(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Tasks completed: %d\n", progress.Data.TasksCompleted)
func (h *Handle) Progress(ctx context.Context) (ProgressResponse, error) {
	resp, err := h.client.V4().GetSherlockProgress(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return ProgressResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSherlockProgressResponse)
	if err != nil {
		return ProgressResponse{ResponseMeta: meta}, err
	}

	return ProgressResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type SherlockTask = v4Client.SherlockTask

type SherlockTasksData = []SherlockTask

type TasksResponse struct {
	Data         SherlockTasksData
	ResponseMeta common.ResponseMeta
}

// Tasks retrieves all tasks for the sherlock.
//
// Example:
//
//	tasks, err := client.Sherlocks.Sherlock(123).Tasks(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sherlock tasks: %d\n", len(tasks.Data))
func (h *Handle) Tasks(ctx context.Context) (TasksResponse, error) {
	resp, err := h.client.V4().GetSherlockTasks(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return TasksResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSherlockTasksResponse)
	if err != nil {
		return TasksResponse{ResponseMeta: meta}, err
	}

	return TasksResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type TaskFlagResponse = v4Client.TaskFlagResponse

type OwnResponse struct {
	Data         TaskFlagResponse
	ResponseMeta common.ResponseMeta
}

// Own submits a flag for a sherlock task.
//
// Example:
//
//	result, err := client.Sherlocks.Sherlock(123).Own(ctx, 456, "HTB{example_flag}")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Flag accepted: %t\n", result.Data.Success)
func (h *Handle) Own(ctx context.Context, taskId int, flag string) (OwnResponse, error) {
	body := v4Client.PostSherlockTasksFlagJSONRequestBody{
		Flag: flag,
	}
	resp, err := h.client.V4().PostSherlockTasksFlag(
		h.client.Limiter().Wrap(ctx),
		h.id,
		taskId,
		body,
	)
	if err != nil {
		return OwnResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostSherlockTasksFlagResponse)
	if err != nil {
		return OwnResponse{ResponseMeta: meta}, err
	}

	return OwnResponse{
		Data:         *parsed.JSON201,
		ResponseMeta: meta,
	}, nil
}

type SherlockDetailData = v4Client.SherlockDetail

type DetailResponse struct {
	Data         SherlockDetailData
	ResponseMeta common.ResponseMeta
}

// Details retrieves detailed sherlock information from the ID-based endpoint.
//
// Example:
//
//	details, err := client.Sherlocks.Sherlock(123).Details(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sherlock detail ID: %d\n", details.Data.Id)
func (h *Handle) Details(ctx context.Context) (DetailResponse, error) {
	resp, err := h.client.V4().GetSherlockInfo(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return DetailResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSherlockInfoResponse)
	if err != nil {
		return DetailResponse{ResponseMeta: meta}, err
	}

	return DetailResponse{
		Data:         parsed.JSON200.Data,
		ResponseMeta: meta,
	}, nil
}

type WriteupData = v4Client.WriteupData

type WriteupResponse struct {
	Data         WriteupData
	ResponseMeta common.ResponseMeta
}

// Writeup retrieves sherlock writeup metadata.
//
// Example:
//
//	writeup, err := client.Sherlocks.Sherlock(123).Writeup(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Writeup official ID: %d\n", writeup.Data.Official.Id)
func (h *Handle) Writeup(ctx context.Context) (WriteupResponse, error) {
	resp, err := h.client.V4().GetSherlockWriteup(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return WriteupResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSherlockWriteupResponse)
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

// WriteupOfficial downloads the sherlock official writeup payload.
//
// Example:
//
//	writeup, err := client.Sherlocks.Sherlock(123).WriteupOfficial(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Writeup bytes: %d\n", len(writeup.Data))
func (h *Handle) WriteupOfficial(ctx context.Context) (WriteupOfficialResponse, error) {
	resp, err := h.client.V4().GetSherlockWriteupOfficial(
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
