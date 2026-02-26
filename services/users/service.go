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

type UserBadges = v4Client.ProfileBadgesItems

type ProfileBadgesResponse struct {
	Data         UserBadges
	ResponseMeta common.ResponseMeta
}

// ProfileBadges retrieves badge information for the user.
// This includes details about the badges the user has earned.
//
// Example:
//
//	badges, err := client.Users.User(12345).ProfileBadges(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("User has %d badges\n", len(badges.Data))
//	profile, err := client.Users.User(12345).ProfileBasic(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, badge := range badges.Data {
//		fmt.Printf("Badge: %s (Date: %s)\n", badge.Name, badge.Pivot.CreatedAt)
//	}
func (h *Handle) ProfileBadges(ctx context.Context) (ProfileBadgesResponse, error) {
	resp, err := h.client.V4().GetUserProfileBadges(
		h.client.Limiter().Wrap(ctx),
		h.id,
		nil,
	)
	if err != nil {
		return ProfileBadgesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserProfileBadgesResponse)
	if err != nil {
		return ProfileBadgesResponse{ResponseMeta: meta}, err
	}

	return ProfileBadgesResponse{
		Data:         parsed.JSON200.Badges,
		ResponseMeta: meta,
	}, nil
}

// GraphPeriod identifies allowed periods for profile graph requests.
type GraphPeriod = v4Client.GetUserProfileGraphParamsPeriod

const (
	GraphPeriod1W GraphPeriod = v4Client.GetUserProfileGraphParamsPeriodN1W
	GraphPeriod1M GraphPeriod = v4Client.GetUserProfileGraphParamsPeriodN1M
	GraphPeriod3M GraphPeriod = v4Client.GetUserProfileGraphParamsPeriodN3M
	GraphPeriod6M GraphPeriod = v4Client.GetUserProfileGraphParamsPeriodN6M
	GraphPeriod1Y GraphPeriod = v4Client.GetUserProfileGraphParamsPeriodN1Y
)

// ContentType identifies v5 profile content filter values.
type ContentType = v5Client.GetUserProfileContentParamsType

const (
	ContentTypeChallenge ContentType = v5Client.GetUserProfileContentParamsTypeChallenge
	ContentTypeMachine   ContentType = v5Client.GetUserProfileContentParamsTypeMachine
	ContentTypeSherlock  ContentType = v5Client.GetUserProfileContentParamsTypeSherlock
)

type UserAchievementData = v4Client.UserAchievementTarTypeUserIdTarIdResponse

// AchievementResponse contains achievement information for a user and target.
type AchievementResponse struct {
	Data         UserAchievementData
	ResponseMeta common.ResponseMeta
}

// Achievement retrieves achievement information for a user target tuple.
//
// Example:
//
//	achievement, err := client.Users.User(12345).Achievement(ctx, "machine", 54321)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Achievement data: %+v\n", achievement.Data)
func (h *Handle) Achievement(ctx context.Context, targetType string, targetId int) (AchievementResponse, error) {
	resp, err := h.client.V4().GetUserAchievement(
		h.client.Limiter().Wrap(ctx),
		targetType,
		h.id,
		targetId,
	)
	if err != nil {
		return AchievementResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserAchievementResponse)
	if err != nil {
		return AchievementResponse{ResponseMeta: meta}, err
	}

	return AchievementResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ConnectionStatusData = v4Client.UserConnectionStatusResponse

// ConnectionStatusResponse contains VPN connection status for the current user.
type ConnectionStatusResponse struct {
	Data         ConnectionStatusData
	ResponseMeta common.ResponseMeta
}

// ConnectionStatus retrieves current user's connection status.
//
// Example:
//
//	status, err := client.Users.ConnectionStatus(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Connection status: %+v\n", status.Data)
func (s *Service) ConnectionStatus(ctx context.Context) (ConnectionStatusResponse, error) {
	resp, err := s.base.Client.V4().GetUserConnectionStatus(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ConnectionStatusResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserConnectionStatusResponse)
	if err != nil {
		return ConnectionStatusResponse{ResponseMeta: meta}, err
	}

	return ConnectionStatusResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type DashboardData = v4Client.UserDashboardResponse

// DashboardResponse contains dashboard payload for current user.
type DashboardResponse struct {
	Data         DashboardData
	ResponseMeta common.ResponseMeta
}

// Dashboard retrieves v4 user dashboard data.
//
// Example:
//
//	dashboard, err := client.Users.Dashboard(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Dashboard data: %+v\n", dashboard.Data)
func (s *Service) Dashboard(ctx context.Context) (DashboardResponse, error) {
	resp, err := s.base.Client.V4().GetUserDashboard(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return DashboardResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserDashboardResponse)
	if err != nil {
		return DashboardResponse{ResponseMeta: meta}, err
	}

	return DashboardResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type DashboardTabloidData = v4Client.UserDashboardTabloidResponse

// DashboardTabloidResponse contains dashboard tabloid payload.
type DashboardTabloidResponse struct {
	Data         DashboardTabloidData
	ResponseMeta common.ResponseMeta
}

// DashboardTabloid retrieves v4 user dashboard tabloid data.
//
// Example:
//
//	tabloid, err := client.Users.DashboardTabloid(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Dashboard tabloid: %+v\n", tabloid.Data)
func (s *Service) DashboardTabloid(ctx context.Context) (DashboardTabloidResponse, error) {
	resp, err := s.base.Client.V4().GetUserDashboardTabloid(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return DashboardTabloidResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserDashboardTabloidResponse)
	if err != nil {
		return DashboardTabloidResponse{ResponseMeta: meta}, err
	}

	return DashboardTabloidResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type DashboardFavoritesData = v5Client.UserDashboardFavoritesResponse

// DashboardFavoritesResponse contains v5 dashboard favorites payload.
type DashboardFavoritesResponse struct {
	Data         DashboardFavoritesData
	ResponseMeta common.ResponseMeta
}

// DashboardFavorites retrieves v5 dashboard favorites.
//
// Example:
//
//	favorites, err := client.Users.DashboardFavorites(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Dashboard favorites: %+v\n", favorites.Data)
func (s *Service) DashboardFavorites(ctx context.Context) (DashboardFavoritesResponse, error) {
	resp, err := s.base.Client.V5().GetUserDashboardFavorites(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return DashboardFavoritesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v5Client.ParseGetUserDashboardFavoritesResponse)
	if err != nil {
		return DashboardFavoritesResponse{ResponseMeta: meta}, err
	}

	return DashboardFavoritesResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type DashboardInProgressData = v5Client.UserDashboardInProgressResponse

// DashboardInProgressResponse contains v5 dashboard in-progress payload.
type DashboardInProgressResponse struct {
	Data         DashboardInProgressData
	ResponseMeta common.ResponseMeta
}

// DashboardInProgress retrieves v5 dashboard in-progress data.
//
// Example:
//
//	inProgress, err := client.Users.DashboardInProgress(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Dashboard in progress: %+v\n", inProgress.Data)
func (s *Service) DashboardInProgress(ctx context.Context) (DashboardInProgressResponse, error) {
	resp, err := s.base.Client.V5().GetUserDashboardInProgress(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return DashboardInProgressResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v5Client.ParseGetUserDashboardInProgressResponse)
	if err != nil {
		return DashboardInProgressResponse{ResponseMeta: meta}, err
	}

	return DashboardInProgressResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type DashboardRecommendedData = v5Client.UserDashboardRecommendedResponse

// DashboardRecommendedResponse contains v5 dashboard recommended payload.
type DashboardRecommendedResponse struct {
	Data         DashboardRecommendedData
	ResponseMeta common.ResponseMeta
}

// DashboardRecommended retrieves v5 dashboard recommended data.
//
// Example:
//
//	recommended, err := client.Users.DashboardRecommended(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Dashboard recommended: %+v\n", recommended.Data)
func (s *Service) DashboardRecommended(ctx context.Context) (DashboardRecommendedResponse, error) {
	resp, err := s.base.Client.V5().GetUserDashboardRecommended(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return DashboardRecommendedResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v5Client.ParseGetUserDashboardRecommendedResponse)
	if err != nil {
		return DashboardRecommendedResponse{ResponseMeta: meta}, err
	}

	return DashboardRecommendedResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

// Disrespect removes respect from the selected user.
//
// Example:
//
//	result, err := client.Users.User(12345).Disrespect(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Disrespect result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (h *Handle) Disrespect(ctx context.Context) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostUserDisrespect(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return common.MessageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostUserDisrespectResponse)
	if err != nil {
		return common.MessageResponse{ResponseMeta: meta}, err
	}

	return common.MessageResponse{
		Data:         common.Message{Message: parsed.JSON200.Message, Success: parsed.JSON200.Success},
		ResponseMeta: meta,
	}, nil
}

// Follow follows the selected user.
//
// Example:
//
//	result, err := client.Users.User(12345).Follow(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Follow result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (h *Handle) Follow(ctx context.Context) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostUserFollow(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return common.MessageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostUserFollowResponse)
	if err != nil {
		return common.MessageResponse{ResponseMeta: meta}, err
	}

	return common.MessageResponse{
		Data:         common.Message{Message: parsed.JSON200.Message, Success: parsed.JSON200.Success},
		ResponseMeta: meta,
	}, nil
}

type FollowersData = v4Client.UserFollowersResponse

// FollowersResponse contains current user's follower list payload.
type FollowersResponse struct {
	Data         FollowersData
	ResponseMeta common.ResponseMeta
}

// Followers retrieves current user's followers.
//
// Example:
//
//	followers, err := client.Users.Followers(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Followers payload: %+v\n", followers.Data)
func (s *Service) Followers(ctx context.Context) (FollowersResponse, error) {
	resp, err := s.base.Client.V4().GetUserFollowers(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return FollowersResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserFollowersResponse)
	if err != nil {
		return FollowersResponse{ResponseMeta: meta}, err
	}

	return FollowersResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type InfoData = v4Client.UserInfoResponse

// InfoResponse contains current authenticated user info.
type InfoResponse struct {
	Data         InfoData
	ResponseMeta common.ResponseMeta
}

// Info retrieves current authenticated user info.
//
// Example:
//
//	info, err := client.Users.Info(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("User info: %+v\n", info.Data)
func (s *Service) Info(ctx context.Context) (InfoResponse, error) {
	resp, err := s.base.Client.V4().GetUserInfo(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return InfoResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserInfoResponse)
	if err != nil {
		return InfoResponse{ResponseMeta: meta}, err
	}

	return InfoResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ProfileBloodsData = v4Client.ProfileBloodsUserIdResponse

// ProfileBloodsResponse contains first-blood statistics for a user.
type ProfileBloodsResponse struct {
	Data         ProfileBloodsData
	ResponseMeta common.ResponseMeta
}

// ProfileBloods retrieves blood information for the selected user.
//
// Example:
//
//	bloods, err := client.Users.User(12345).ProfileBloods(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Profile bloods: %+v\n", bloods.Data)
func (h *Handle) ProfileBloods(ctx context.Context) (ProfileBloodsResponse, error) {
	resp, err := h.client.V4().GetUserProfileBloods(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return ProfileBloodsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserProfileBloodsResponse)
	if err != nil {
		return ProfileBloodsResponse{ResponseMeta: meta}, err
	}

	return ProfileBloodsResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ProfileChartMachinesAttackData = v4Client.ProfileChartMachinesAttackUserIdResponse

// ProfileChartMachinesAttackResponse contains machine attack chart for a user.
type ProfileChartMachinesAttackResponse struct {
	Data         ProfileChartMachinesAttackData
	ResponseMeta common.ResponseMeta
}

// ProfileChartMachinesAttack retrieves machine attack chart data for the selected user.
//
// Example:
//
//	chart, err := client.Users.User(12345).ProfileChartMachinesAttack(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Machine attack chart: %+v\n", chart.Data)
func (h *Handle) ProfileChartMachinesAttack(ctx context.Context) (ProfileChartMachinesAttackResponse, error) {
	resp, err := h.client.V4().GetUserProfileChartMachinesAttack(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return ProfileChartMachinesAttackResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserProfileChartMachinesAttackResponse)
	if err != nil {
		return ProfileChartMachinesAttackResponse{ResponseMeta: meta}, err
	}

	return ProfileChartMachinesAttackResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ProfileContentData = v5Client.UserProfileContentResponse

// ProfileContentResponse contains v5 profile content payload.
type ProfileContentResponse struct {
	Data         ProfileContentData
	ResponseMeta common.ResponseMeta
}

// ProfileContent retrieves v5 content items for the selected user.
//
// Example:
//
//	content, err := client.Users.User(12345).ProfileContent(ctx, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Profile content: %+v\n", content.Data)
func (h *Handle) ProfileContent(ctx context.Context, params *v5Client.GetUserProfileContentParams) (ProfileContentResponse, error) {
	resp, err := h.client.V5().GetUserProfileContent(h.client.Limiter().Wrap(ctx), h.id, params)
	if err != nil {
		return ProfileContentResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v5Client.ParseGetUserProfileContentResponse)
	if err != nil {
		return ProfileContentResponse{ResponseMeta: meta}, err
	}

	return ProfileContentResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ProfileGraphData = v4Client.ProfileGraphPeriodUserIdResponse

// ProfileGraphResponse contains profile graph data for the selected user.
type ProfileGraphResponse struct {
	Data         ProfileGraphData
	ResponseMeta common.ResponseMeta
}

// ProfileGraph retrieves profile graph data for the selected user and period.
//
// Example:
//
//	graph, err := client.Users.User(12345).ProfileGraph(ctx, users.GraphPeriod1Y)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Profile graph: %+v\n", graph.Data)
func (h *Handle) ProfileGraph(ctx context.Context, period GraphPeriod) (ProfileGraphResponse, error) {
	resp, err := h.client.V4().GetUserProfileGraph(h.client.Limiter().Wrap(ctx), v4Client.GetUserProfileGraphParamsPeriod(period), h.id)
	if err != nil {
		return ProfileGraphResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserProfileGraphResponse)
	if err != nil {
		return ProfileGraphResponse{ResponseMeta: meta}, err
	}

	return ProfileGraphResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ProfileProgressChallengesData = v4Client.ProfileProgressChallengesUserIdResponse

// ProfileProgressChallengesResponse contains challenge progress data.
type ProfileProgressChallengesResponse struct {
	Data         ProfileProgressChallengesData
	ResponseMeta common.ResponseMeta
}

// ProfileProgressChallenges retrieves challenge progress for the selected user.
//
// Example:
//
//	progress, err := client.Users.User(12345).ProfileProgressChallenges(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Challenge progress: %+v\n", progress.Data)
func (h *Handle) ProfileProgressChallenges(ctx context.Context) (ProfileProgressChallengesResponse, error) {
	resp, err := h.client.V4().GetUserProfileProgressChallenges(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return ProfileProgressChallengesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserProfileProgressChallengesResponse)
	if err != nil {
		return ProfileProgressChallengesResponse{ResponseMeta: meta}, err
	}

	return ProfileProgressChallengesResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ProfileProgressFortressData = v4Client.ProfileProgressFortressUserIdResponse

// ProfileProgressFortressResponse contains fortress progress data.
type ProfileProgressFortressResponse struct {
	Data         ProfileProgressFortressData
	ResponseMeta common.ResponseMeta
}

// ProfileProgressFortress retrieves fortress progress for the selected user.
//
// Example:
//
//	progress, err := client.Users.User(12345).ProfileProgressFortress(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Fortress progress: %+v\n", progress.Data)
func (h *Handle) ProfileProgressFortress(ctx context.Context) (ProfileProgressFortressResponse, error) {
	resp, err := h.client.V4().GetUserProfileProgressFortress(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return ProfileProgressFortressResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserProfileProgressFortressResponse)
	if err != nil {
		return ProfileProgressFortressResponse{ResponseMeta: meta}, err
	}

	return ProfileProgressFortressResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ProfileProgressProlabData = v4Client.ProfileProgressProlabUserIdResponse

// ProfileProgressProlabResponse contains prolab progress data.
type ProfileProgressProlabResponse struct {
	Data         ProfileProgressProlabData
	ResponseMeta common.ResponseMeta
}

// ProfileProgressProlab retrieves prolab progress for the selected user.
//
// Example:
//
//	progress, err := client.Users.User(12345).ProfileProgressProlab(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Prolab progress: %+v\n", progress.Data)
func (h *Handle) ProfileProgressProlab(ctx context.Context) (ProfileProgressProlabResponse, error) {
	resp, err := h.client.V4().GetUserProfileProgressProlab(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return ProfileProgressProlabResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserProfileProgressProlabResponse)
	if err != nil {
		return ProfileProgressProlabResponse{ResponseMeta: meta}, err
	}

	return ProfileProgressProlabResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ProfileProgressSherlocksData = v4Client.UserProfileProgressSherlocksResponse

// ProfileProgressSherlocksResponse contains sherlocks progress data.
type ProfileProgressSherlocksResponse struct {
	Data         ProfileProgressSherlocksData
	ResponseMeta common.ResponseMeta
}

// ProfileProgressSherlocks retrieves sherlocks progress for the selected user.
//
// Example:
//
//	progress, err := client.Users.User(12345).ProfileProgressSherlocks(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sherlocks progress: %+v\n", progress.Data)
func (h *Handle) ProfileProgressSherlocks(ctx context.Context) (ProfileProgressSherlocksResponse, error) {
	resp, err := h.client.V4().GetUserProfileProgressSherlocks(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return ProfileProgressSherlocksResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserProfileProgressSherlocksResponse)
	if err != nil {
		return ProfileProgressSherlocksResponse{ResponseMeta: meta}, err
	}

	return ProfileProgressSherlocksResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ProfileSummaryData = v4Client.UserProfileSummaryResponse

// ProfileSummaryResponse contains compact profile summary for the current authenticated user.
type ProfileSummaryResponse struct {
	Data         ProfileSummaryData
	ResponseMeta common.ResponseMeta
}

// ProfileSummary retrieves profile summary for the current authenticated user.
//
// Example:
//
//	summary, err := client.Users.ProfileSummary(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Profile summary: %+v\n", summary.Data)
func (s *Service) ProfileSummary(ctx context.Context) (ProfileSummaryResponse, error) {
	resp, err := s.base.Client.V4().GetUserProfileSummary(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ProfileSummaryResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserProfileSummaryResponse)
	if err != nil {
		return ProfileSummaryResponse{ResponseMeta: meta}, err
	}

	return ProfileSummaryResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

// Respect gives respect to the selected user.
//
// Example:
//
//	result, err := client.Users.User(12345).Respect(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Respect result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (h *Handle) Respect(ctx context.Context) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostUserRespect(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return common.MessageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostUserRespectResponse)
	if err != nil {
		return common.MessageResponse{ResponseMeta: meta}, err
	}

	return common.MessageResponse{
		Data:         common.Message{Message: parsed.JSON200.Message, Success: parsed.JSON200.Success},
		ResponseMeta: meta,
	}, nil
}

type SettingsData = v4Client.UserSettingsResponse

// SettingsResponse contains current user settings payload.
type SettingsResponse struct {
	Data         SettingsData
	ResponseMeta common.ResponseMeta
}

// Settings retrieves current user settings.
//
// Example:
//
//	settings, err := client.Users.Settings(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Settings payload: %+v\n", settings.Data)
func (s *Service) Settings(ctx context.Context) (SettingsResponse, error) {
	resp, err := s.base.Client.V4().GetUserSettings(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return SettingsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserSettingsResponse)
	if err != nil {
		return SettingsResponse{ResponseMeta: meta}, err
	}

	return SettingsResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type TracksData = v4Client.UserTracksResponse

// TracksResponse contains current user track progress payload.
type TracksResponse struct {
	Data         TracksData
	ResponseMeta common.ResponseMeta
}

// Tracks retrieves current user tracks.
//
// Example:
//
//	tracks, err := client.Users.Tracks(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Tracks payload: %+v\n", tracks.Data)
func (s *Service) Tracks(ctx context.Context) (TracksResponse, error) {
	resp, err := s.base.Client.V4().GetUserTracks(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return TracksResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetUserTracksResponse)
	if err != nil {
		return TracksResponse{ResponseMeta: meta}, err
	}

	return TracksResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

// Unfollow unfollows the selected user.
//
// Example:
//
//	result, err := client.Users.User(12345).Unfollow(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Unfollow result: %s (Success: %t)\n", result.Data.Message, result.Data.Success)
func (h *Handle) Unfollow(ctx context.Context) (common.MessageResponse, error) {
	resp, err := h.client.V4().PostUserUnfollow(h.client.Limiter().Wrap(ctx), h.id)
	if err != nil {
		return common.MessageResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParsePostUserUnfollowResponse)
	if err != nil {
		return common.MessageResponse{ResponseMeta: meta}, err
	}

	return common.MessageResponse{
		Data:         common.Message{Message: parsed.JSON200.Message, Success: parsed.JSON200.Success},
		ResponseMeta: meta,
	}, nil
}
