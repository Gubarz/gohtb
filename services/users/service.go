package users

import (
	"context"

	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/convert"
	"github.com/gubarz/gohtb/internal/errutil"
	"github.com/gubarz/gohtb/internal/extract"
	"github.com/gubarz/gohtb/internal/service"
)

func NewService(client service.Client) *Service {
	return &Service{
		base: service.NewBase(client),
	}
}

// User returns a handle for a specific user with the given ID.
// This handle can be used to perform operations related to that user,
// such as retrieving profile information and activity data.
func (s *Service) User(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

// ProfileBasic retrieves basic profile information for the user.
// This includes fundamental user details such as username, avatar,
// rank, and other core profile data.
//
// Example:
//
//	profile, err := client.Users.User(12345).ProfileBasic(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("User: %s (Rank: %d)\n", profile.Data.Username, profile.Data.Rank)
func (h *Handle) ProfileBasic(ctx context.Context) (ProfileBasicResponse, error) {
	resp, err := h.client.V4().GetUserProfileBasicWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)

	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) ProfileBasicResponse {
			return ProfileBasicResponse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return ProfileBasicResponse{
		Data: fromAPIUserProfile(resp.JSON200.Profile),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}

// ProfileActivity retrieves the activity history for the user.
// This includes recent actions, submissions, and other user activities
// on the HackTheBox platform.
//
// Example:
//
//	activity, err := client.Users.User(12345).ProfileActivity(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, act := range activity.Data {
//		fmt.Printf("Activity: %s at %s\n", act.Type, act.Date)
//	}
func (h *Handle) ProfileActivity(ctx context.Context) (ProfileActivityResposnse, error) {
	resp, err := h.client.V4().GetUserProfileActivityWithResponse(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	raw := extract.Raw(resp)

	if err != nil || resp == nil || resp.JSON200 == nil {
		return errutil.UnwrapFailure(err, raw, common.SafeStatus(resp), func(raw []byte) ProfileActivityResposnse {
			return ProfileActivityResposnse{ResponseMeta: common.ResponseMeta{Raw: raw}}
		})
	}

	return ProfileActivityResposnse{
		Data: convert.SlicePointer(resp.JSON200.Profile.Activity, fromAPIUserActivity),
		ResponseMeta: common.ResponseMeta{
			Raw:        raw,
			StatusCode: resp.StatusCode(),
			Headers:    resp.HTTPResponse.Header,
		},
	}, nil
}
