package users

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	v5Client "github.com/gubarz/gohtb/httpclient/v5"
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

// User returns a handle for a specific user with the given ID.
// This handle can be used to perform operations related to that user,
// such as retrieving profile information and activity data.
func (s *Service) User(id int) *Handle {
	return &Handle{
		client: s.base.Client,
		id:     id,
	}
}

type UserProfile = v4Client.UserProfile

type ProfileBasicResponse struct {
	Data         UserProfile
	ResponseMeta common.ResponseMeta
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
	resp, err := h.client.V4().GetUserProfileBasic(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return ProfileBasicResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserProfileBasicResponse)
	if err != nil {
		return ProfileBasicResponse{ResponseMeta: meta}, err
	}

	return ProfileBasicResponse{
		Data:         parsed.JSON200.Profile,
		ResponseMeta: meta,
	}, nil
}

type UserActivityItem = v4Client.UserActivityItem

type ProfileActivityResponse struct {
	Data         []UserActivityItem
	ResponseMeta common.ResponseMeta
}

// ProfileActivityDeprecated retrieves the activity history for the user.
// This includes recent actions, submissions, and other user activities
// on the HackTheBox platform.
//
// Note: This method is considered deprecated. If you want starting point
// activities you still need to use this method.
//
// Example:
//
//	activity, err := client.Users.User(12345).ProfileActivityDeprecated(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, act := range activity.Data {
//		fmt.Printf("Activity: %s at %s\n", act.Type, act.Date)
//	}
func (h *Handle) ProfileActivityDeprecated(ctx context.Context) (ProfileActivityResponse, error) {
	resp, err := h.client.V4().GetUserProfileActivity(
		h.client.Limiter().Wrap(ctx),
		h.id,
	)
	if err != nil {
		return ProfileActivityResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserProfileActivityResponse)
	if err != nil {
		return ProfileActivityResponse{ResponseMeta: meta}, err
	}

	return ProfileActivityResponse{
		Data:         parsed.JSON200.Profile.Activity,
		ResponseMeta: meta,
	}, nil
}

type UserProfileActivityItem = v5Client.UserProfileActivityItem

type UserProfileActivity struct {
	v5Client.UserProfileActivityBase
	Challenge  v5Client.UserProfileActivityChallenge
	Fortress   v5Client.UserProfileActivityFortress
	MachineOwn v5Client.UserProfileActivityMachineOwn
	Prolab     v5Client.UserProfileActivityProlab
	Sherlock   v5Client.UserProfileActivitySherlock
}

func (a UserProfileActivity) AsChallenge() (v5Client.UserProfileActivityChallenge, bool) {
	return a.Challenge, a.Type == string(v5Client.UserProfileActivityChallengeTypeChallenge)
}

func (a UserProfileActivity) AsFortress() (v5Client.UserProfileActivityFortress, bool) {
	return a.Fortress, a.Type == "fortress"
}

func (a UserProfileActivity) AsMachineOwn() (v5Client.UserProfileActivityMachineOwn, bool) {
	switch a.Type {
	case string(v5Client.UserProfileActivityMachineOwnTypeRoot), string(v5Client.UserProfileActivityMachineOwnTypeUser):
		return a.MachineOwn, true
	default:
		return v5Client.UserProfileActivityMachineOwn{}, false
	}
}

func (a UserProfileActivity) AsProlab() (v5Client.UserProfileActivityProlab, bool) {
	return a.Prolab, a.Type == string(v5Client.UserProfileActivityProlabTypeProlab)
}

func (a UserProfileActivity) AsSherlock() (v5Client.UserProfileActivitySherlock, bool) {
	return a.Sherlock, a.Type == string(v5Client.UserProfileActivitySherlockTypeSherlock)
}

func (h *Handle) ProfileActivity() *UserProfileActivityQuery {
	return &UserProfileActivityQuery{
		client:  h.client,
		id:      h.id,
		page:    1,
		perPage: 100,
	}
}
