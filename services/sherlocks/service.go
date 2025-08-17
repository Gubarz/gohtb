package sherlocks

import (
	"context"
	"strconv"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/service"
)

func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

func (s *Service) Sherlock(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
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
		Data:         fromAPISherlock(parsed.JSON200.Data),
		ResponseMeta: meta,
	}, nil
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
		Data:         fromAPISherlockDownloadLink(parsed.JSON200),
		ResponseMeta: meta,
	}, nil
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
		Data:         fromAPISherlockProgress(parsed.JSON200.Data),
		ResponseMeta: meta,
	}, nil
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
		Data:         convert.Slice(*parsed.JSON200.Data, fromAPISherlockTasks),
		ResponseMeta: meta,
	}, nil
}

func (h *Handle) Own(ctx context.Context, taskId int, flag string) (OwnResponse, error) {
	body := v4Client.PostSherlockTasksFlagJSONRequestBody{
		Flag: &flag,
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
		Data:         fromAPIOwnTask(parsed.JSON201),
		ResponseMeta: meta,
	}, nil
}
