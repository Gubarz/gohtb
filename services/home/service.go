package home

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

// Service provides access to home dashboard endpoints.
type Service struct {
	base service.Base
}

// NewService creates a new home service bound to a shared client.
//
// Example:
//
//	homeService := home.NewService(client)
//	_ = homeService
func NewService(client service.Client) *Service {
	return &Service{base: service.NewBase(client)}
}

type BannerData = v4Client.HomeBannersResponse

// BannerResponse contains home banner payload.
type BannerResponse struct {
	Data         BannerData
	ResponseMeta common.ResponseMeta
}

// Banner retrieves home banner content.
//
// Example:
//
//	banner, err := client.Home.Banner(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Home banner payload: %+v\n", banner.Data)
func (s *Service) Banner(ctx context.Context) (BannerResponse, error) {
	resp, err := s.base.Client.V4().GetHomeBanner(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return BannerResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetHomeBannerResponse)
	if err != nil {
		return BannerResponse{ResponseMeta: meta}, err
	}

	return BannerResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type RecommendedData = v4Client.HomeRecommendedResponse

// RecommendedResponse contains home recommended cards.
type RecommendedResponse struct {
	Data         RecommendedData
	ResponseMeta common.ResponseMeta
}

// Recommended retrieves recommended content shown on the home dashboard.
//
// Example:
//
//	recommended, err := client.Home.Recommended(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Home recommended payload: %+v\n", recommended.Data)
func (s *Service) Recommended(ctx context.Context) (RecommendedResponse, error) {
	resp, err := s.base.Client.V4().GetHomeRecommended(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return RecommendedResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetHomeRecommendedResponse)
	if err != nil {
		return RecommendedResponse{ResponseMeta: meta}, err
	}

	return RecommendedResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type UserProgressData = v4Client.HomeUserProgressResponse

// UserProgressResponse contains user home progress cards.
type UserProgressResponse struct {
	Data         UserProgressData
	ResponseMeta common.ResponseMeta
}

// UserProgress retrieves authenticated user progress cards shown on home.
//
// Example:
//
//	progress, err := client.Home.UserProgress(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Home user progress payload: %+v\n", progress.Data)
func (s *Service) UserProgress(ctx context.Context) (UserProgressResponse, error) {
	resp, err := s.base.Client.V4().GetHomeUserProgress(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return UserProgressResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetHomeUserProgressResponse)
	if err != nil {
		return UserProgressResponse{ResponseMeta: meta}, err
	}

	return UserProgressResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type UserTodoData = v4Client.HomeUserTodoRepsonse

// UserTodoResponse contains user home todo cards.
type UserTodoResponse struct {
	Data         UserTodoData
	ResponseMeta common.ResponseMeta
}

// UserTodo retrieves authenticated user todo cards shown on home.
//
// Example:
//
//	todo, err := client.Home.UserToDo(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Home user todo payload: %+v\n", todo.Data)
func (s *Service) UserToDo(ctx context.Context) (UserTodoResponse, error) {
	resp, err := s.base.Client.V4().GetHomeUserTodo(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return UserTodoResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetHomeUserTodoResponse)
	if err != nil {
		return UserTodoResponse{ResponseMeta: meta}, err
	}

	return UserTodoResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}
