package platform

import (
	"context"

	v4Client "github.com/gubarz/gohtb/httpclient/v4"
	"github.com/gubarz/gohtb/internal/common"
	"github.com/gubarz/gohtb/internal/service"
)

// Service provides access to general platform endpoints.
type Service struct {
	base service.Base
}

// NewService creates a new platform service bound to a shared client.
//
// Example:
//
//	platformService := platform.NewService(client)
//	_ = platformService
func NewService(client service.Client) *Service {
	return &Service{base: service.NewBase(client)}
}

type AnnouncementsData = v4Client.AnnouncementResponse

// AnnouncementsResponse contains platform announcement data.
type AnnouncementsResponse struct {
	Data         AnnouncementsData
	ResponseMeta common.ResponseMeta
}

// Announcements retrieves public platform announcements.
//
// Example:
//
//	announcements, err := client.Platform.Announcements(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Announcements payload: %+v\n", announcements.Data)
func (s *Service) Announcements(ctx context.Context) (AnnouncementsResponse, error) {
	resp, err := s.base.Client.V4().GetAnnouncements(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return AnnouncementsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetAnnouncementsResponse)
	if err != nil {
		return AnnouncementsResponse{ResponseMeta: meta}, err
	}

	return AnnouncementsResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ChangelogsData = v4Client.ChangelogsResponse

// ChangelogsResponse contains global changelog entries.
type ChangelogsResponse struct {
	Data         ChangelogsData
	ResponseMeta common.ResponseMeta
}

// Changelogs retrieves global platform changelog entries.
//
// Example:
//
//	changelogs, err := client.Platform.Changelogs(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Changelogs payload: %+v\n", changelogs.Data)
func (s *Service) Changelogs(ctx context.Context) (ChangelogsResponse, error) {
	resp, err := s.base.Client.V4().GetChangelogs(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ChangelogsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetChangelogsResponse)
	if err != nil {
		return ChangelogsResponse{ResponseMeta: meta}, err
	}

	return ChangelogsResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type ContentStatsData = v4Client.ContentStatsResponse

// ContentStatsResponse contains content statistics.
type ContentStatsResponse struct {
	Data         ContentStatsData
	ResponseMeta common.ResponseMeta
}

// ContentStats retrieves high-level platform content counts.
//
// Example:
//
//	stats, err := client.Platform.ContentStats(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Content stats: %+v\n", stats.Data)
func (s *Service) ContentStats(ctx context.Context) (ContentStatsResponse, error) {
	resp, err := s.base.Client.V4().GetContentStats(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return ContentStatsResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetContentStatsResponse)
	if err != nil {
		return ContentStatsResponse{ResponseMeta: meta}, err
	}

	return ContentStatsResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type NavigationMainData = v4Client.NavigationMainResponse

// NavigationMainResponse contains navigation summary data.
type NavigationMainResponse struct {
	Data         NavigationMainData
	ResponseMeta common.ResponseMeta
}

// NavigationMain retrieves main navigation summary information.
//
// Example:
//
//	navigation, err := client.Platform.NavigationMain(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Navigation payload: %+v\n", navigation.Data)
func (s *Service) NavigationMain(ctx context.Context) (NavigationMainResponse, error) {
	resp, err := s.base.Client.V4().GetNavigationMain(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return NavigationMainResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetNavigationMainResponse)
	if err != nil {
		return NavigationMainResponse{ResponseMeta: meta}, err
	}

	return NavigationMainResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type NoticesData = v4Client.NoticesResponse

// NoticesResponse contains active platform notices.
type NoticesResponse struct {
	Data         NoticesData
	ResponseMeta common.ResponseMeta
}

// Notices retrieves active platform notices.
//
// Example:
//
//	notices, err := client.Platform.Notices(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Notices payload: %+v\n", notices.Data)
func (s *Service) Notices(ctx context.Context) (NoticesResponse, error) {
	resp, err := s.base.Client.V4().GetNotices(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return NoticesResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetNoticesResponse)
	if err != nil {
		return NoticesResponse{ResponseMeta: meta}, err
	}

	return NoticesResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}

type SidebarAnnouncementData = v4Client.SidebarAnnouncementResponse

// SidebarAnnouncementResponse contains sidebar announcement data.
type SidebarAnnouncementResponse struct {
	Data         SidebarAnnouncementData
	ResponseMeta common.ResponseMeta
}

// SidebarAnnouncement retrieves the sidebar announcement payload.
//
// Example:
//
//	sidebar, err := client.Platform.SidebarAnnouncement(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Sidebar announcement payload: %+v\n", sidebar.Data)
func (s *Service) SidebarAnnouncement(ctx context.Context) (SidebarAnnouncementResponse, error) {
	resp, err := s.base.Client.V4().GetSidebarAnnouncement(s.base.Client.Limiter().Wrap(ctx))
	if err != nil {
		return SidebarAnnouncementResponse{ResponseMeta: common.ResponseMeta{}}, err
	}

	parsed, meta, err := common.Parse(resp, v4Client.ParseGetSidebarAnnouncementResponse)
	if err != nil {
		return SidebarAnnouncementResponse{ResponseMeta: meta}, err
	}

	return SidebarAnnouncementResponse{Data: *parsed.JSON200, ResponseMeta: meta}, nil
}
