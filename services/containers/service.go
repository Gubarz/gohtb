package containers

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

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

// Container returns a handle for a specific container with the given ID.
// This handle can be used to perform operations on the container such as
// retrieving information, starting/stopping instances, or submitting flags.
func (s *Service) Container(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

// Start initiates a new instance of the container.
// This creates and starts a container instance for solving.
//
// Example:
//
//	result, err := client.Containers.Container(12345).Start(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Container started: %s\n", result.Data.Message)
func (h *Handle) Start(ctx context.Context) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostContainerStartWithFormdataBody(
		h.client.Limiter().Wrap(ctx),
		v4Client.PostContainerStartFormdataRequestBody{
			ContainerableId: h.id,
		},
	)
	if err != nil {
		return common.MessageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostContainerStartResponse)
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

// Stop terminates the running container instance.
// This stops and destroys the active container instance.
//
// Example:
//
//	result, err := client.Containers.Container(12345).Stop(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Container stopped: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (h *Handle) Stop(ctx context.Context) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostContainerStopWithFormdataBody(
		h.client.Limiter().Wrap(ctx),
		v4Client.PostContainerStopFormdataRequestBody{
			ContainerableId: h.id,
		},
	)
	if err != nil {
		return common.MessageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostContainerStopResponse)
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
