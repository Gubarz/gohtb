package sherlocks

import (
	"context"
	"strconv"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
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

func (s *Service) Sherlock(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

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
