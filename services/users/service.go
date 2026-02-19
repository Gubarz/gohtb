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

// NewService creates a new users service bound to a shared client.
//
// Example:
//
//	userService := users.NewService(client)
//	_ = userService
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
//
// Example:
//
//	user := client.Users.User(12345)
//	profile, err := user.ProfileBasic(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Username: %s\n", profile.Data.Username)
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

// AsChallenge returns challenge activity data if this activity is a challenge event.
//
// Example:
//
//	activity, ok := item.AsChallenge()
//	if ok {
//		fmt.Printf("Challenge activity: %s\n", activity.Name)
//	}
func (a UserProfileActivity) AsChallenge() (v5Client.UserProfileActivityChallenge, bool) {
	return a.Challenge, a.Type == string(v5Client.UserProfileActivityChallengeTypeChallenge)
}

// AsFortress returns fortress activity data if this activity is a fortress event.
//
// Example:
//
//	activity, ok := item.AsFortress()
//	if ok {
//		fmt.Printf("Fortress activity: %s\n", activity.Name)
//	}
func (a UserProfileActivity) AsFortress() (v5Client.UserProfileActivityFortress, bool) {
	return a.Fortress, a.Type == "fortress"
}

// AsMachineOwn returns machine own activity data if this activity is a user/root own event.
//
// Example:
//
//	activity, ok := item.AsMachineOwn()
//	if ok {
//		fmt.Printf("Machine own activity: %s\n", activity.Name)
//	}
func (a UserProfileActivity) AsMachineOwn() (v5Client.UserProfileActivityMachineOwn, bool) {
	switch a.Type {
	case string(v5Client.UserProfileActivityMachineOwnTypeRoot), string(v5Client.UserProfileActivityMachineOwnTypeUser):
		return a.MachineOwn, true
	default:
		return v5Client.UserProfileActivityMachineOwn{}, false
	}
}

// AsProlab returns prolab activity data if this activity is a prolab event.
//
// Example:
//
//	activity, ok := item.AsProlab()
//	if ok {
//		fmt.Printf("Prolab activity: %s\n", activity.Name)
//	}
func (a UserProfileActivity) AsProlab() (v5Client.UserProfileActivityProlab, bool) {
	return a.Prolab, a.Type == string(v5Client.UserProfileActivityProlabTypeProlab)
}

// AsSherlock returns sherlock activity data if this activity is a sherlock event.
//
// Example:
//
//	activity, ok := item.AsSherlock()
//	if ok {
//		fmt.Printf("Sherlock activity: %s\n", activity.Name)
//	}
func (a UserProfileActivity) AsSherlock() (v5Client.UserProfileActivitySherlock, bool) {
	return a.Sherlock, a.Type == string(v5Client.UserProfileActivitySherlockTypeSherlock)
}

// ProfileActivity creates a paginated activity query for this user.
//
// Example:
//
//	activity, err := client.Users.User(12345).ProfileActivity().PerPage(25).Results(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Activity items: %d\n", len(activity.Data))
func (h *Handle) ProfileActivity() *UserProfileActivityQuery {
	return &UserProfileActivityQuery{
		client:  h.client,
		id:      h.id,
		page:    1,
		perPage: 100,
	}
}

type AppToken = v4Client.UserApptokenListItem

type AppTokensResponse struct {
	Data         []AppToken
	ResponseMeta common.ResponseMeta
}

// AppTokens lists API app tokens for the authenticated user.
//
// Example:
//
//	tokens, err := client.Users.AppTokens(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Tokens: %d\n", len(tokens.Data))
func (s *Service) AppTokens(ctx context.Context) (AppTokensResponse, error) {
	resp, err := s.base.Client.V4().GetUserApptokenList(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return AppTokensResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserApptokenListResponse)
	if err != nil {
		return AppTokensResponse{ResponseMeta: meta}, err
	}

	return AppTokensResponse{
		Data:         parsed.JSON200.Tokens,
		ResponseMeta: meta,
	}, nil
}

type AppTokenCreateRequest struct {
	Name        string
	ExpireAfter float32
}

type AppTokenCreateData = v4Client.UserApptokenCreateResponse

type CreateAppTokenResponse struct {
	Data         AppTokenCreateData
	ResponseMeta common.ResponseMeta
}

// CreateAppToken creates a new API app token for the authenticated user.
//
// Example:
//
//	token, err := client.Users.CreateAppToken(ctx, users.AppTokenCreateRequest{
//		Name:        "ci-runner",
//		ExpireAfter: 30,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Token name: %s\n", token.Data.Name)
func (s *Service) CreateAppToken(ctx context.Context, req AppTokenCreateRequest) (CreateAppTokenResponse, error) {
	resp, err := s.base.Client.V4().PostUserApptokenCreate(
		s.base.Client.Limiter().Wrap(ctx),
		v4Client.PostUserApptokenCreateJSONRequestBody{
			Name:        req.Name,
			ExpireAfter: req.ExpireAfter,
		},
	)
	if err != nil {
		return CreateAppTokenResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostUserApptokenCreateResponse)
	if err != nil {
		return CreateAppTokenResponse{ResponseMeta: meta}, err
	}

	return CreateAppTokenResponse{
		Data:         *parsed.JSON200,
		ResponseMeta: meta,
	}, nil
}

type AppTokenDeleteRequest struct {
	Name string
}

// DeleteAppToken deletes an API app token by name for the authenticated user.
//
// Example:
//
//	result, err := client.Users.DeleteAppToken(ctx, users.AppTokenDeleteRequest{
//		Name: "ci-runner",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Delete result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (s *Service) DeleteAppToken(ctx context.Context, req AppTokenDeleteRequest) (common.MessageResponse, error) {
	resp, err := s.base.Client.V4().PostUserApptokenDelete(
		s.base.Client.Limiter().Wrap(ctx),
		v4Client.PostUserApptokenDeleteJSONRequestBody{
			Name: req.Name,
		},
	)
	if err != nil {
		return common.MessageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostUserApptokenDeleteResponse)
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
