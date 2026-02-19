package tracks

import (
	"context"
	"strconv"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

type Service struct {
	base service.Base
}

// NewService creates a new tracks service bound to a shared client.
//
// Example:
//
//	trackService := tracks.NewService(client)
//	_ = trackService
func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

type Handle struct {
	client service.Client
	id     int
}

// Track returns a handle for a specific track with the given ID.
//
// Example:
//
//	track := client.Tracks.Track(42)
//	_ = track
func (s *Service) Track(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

type Track = v4Client.TracksItems

type ListResponse struct {
	Data         []Track
	ResponseMeta common.ResponseMeta
}

// List retrieves all tracks available to the authenticated user.
//
// Example:
//
//	tracks, err := client.Tracks.List(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Tracks: %d\n", len(tracks.Data))
func (s *Service) List(ctx context.Context) (ListResponse, error) {
	resp, err := s.base.Client.V4().GetTracks(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ListResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetTracksResponse)
	if err != nil {
		return ListResponse{ResponseMeta: meta}, err
	}

	return ListResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type TrackDetails = v4Client.TracksIdResponse

type DetailsResponse struct {
	Data         TrackDetails
	ResponseMeta common.ResponseMeta
}

// Info retrieves full details for the track handle.
//
// Example:
//
//	details, err := client.Tracks.Track(42).Info(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	track, err := details.Data.AsTrackSuccessResponse()
//	if err == nil {
//		fmt.Printf("Track: %s\n", track.Name)
//	}
func (h *Handle) Info(ctx context.Context) (DetailsResponse, error) {
	resp, err := h.client.V4().GetTracksId(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return DetailsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetTracksIdResponse)
	if err != nil {
		return DetailsResponse{ResponseMeta: meta}, err
	}

	return DetailsResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type EnrollData = v4Client.TracksEnrollResponse

type EnrollResponse struct {
	Data         EnrollData
	ResponseMeta common.ResponseMeta
}

// Enroll enrolls the authenticated user into the track.
//
// Example:
//
//	enroll, err := client.Tracks.Track(42).Enroll(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Enrolled: %t\n", enroll.Data.Enrolled)
func (h *Handle) Enroll(ctx context.Context) (EnrollResponse, error) {
	resp, err := h.client.V4().PostTracksEnroll(
		h.client.Limiter().Wrap(ctx),
		strconv.Itoa(h.id),
	)
	if err != nil {
		return EnrollResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostTracksEnrollResponse)
	if err != nil {
		return EnrollResponse{ResponseMeta: meta}, err
	}

	return EnrollResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type LikeData = v4Client.TracksLikeResponse

type LikeResponse struct {
	Data         LikeData
	ResponseMeta common.ResponseMeta
}

// Like toggles like status for the track for the authenticated user.
//
// Example:
//
//	like, err := client.Tracks.Track(42).Like(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Liked: %t\n", like.Data.Liked)
func (h *Handle) Like(ctx context.Context) (LikeResponse, error) {
	resp, err := h.client.V4().PostTracksLike(
		h.client.Limiter().Wrap(ctx),
		strconv.Itoa(h.id),
	)
	if err != nil {
		return LikeResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostTracksLikeResponse)
	if err != nil {
		return LikeResponse{ResponseMeta: meta}, err
	}

	return LikeResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}
