package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/community"
	"github.com/sarbojitrana/nexus/internal/model/post"
	"github.com/sarbojitrana/nexus/internal/repository"
	"github.com/sarbojitrana/nexus/internal/server"
)

type CommunityService struct {
	server *server.Server
	repo   *repository.CommunityRepository
}

func NewCommunityService(s *server.Server, repo *repository.CommunityRepository) *CommunityService {
	return &CommunityService{
		server: s,
		repo:   repo,
	}
}

func (s *CommunityService) CreateCommunity(ctx context.Context, adminID string, payload *community.CreateCommunityPayload) (*community.Community, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	payload.AdminID = adminID // never trust a client-supplied admin id

	created, err := s.repo.CreateCommunity(ctx, payload)
	if err != nil {
		logger.Error().Err(err).Str("admin_id", adminID).Str("slug", payload.Slug).Msg("failed to create community")
		return nil, err
	}

	logger.Info().Str("event", "community_created").Str("community_id", created.ID.String()).Str("admin_id", adminID).Msg("community created")
	return created, nil
}

func (s *CommunityService) GetCommunities(ctx context.Context, query *community.GetCommunitiesQuery) (*model.CursorPaginatedResponse[community.MiniCommunity], error) {
	res, err := s.repo.GetCommunities(ctx, query)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Msg("failed to list communities")
		return nil, err
	}
	return res, nil
}

// GetByID composes a CommunityResponse with the viewer's role, when known.
func (s *CommunityService) GetByID(ctx context.Context, viewerID *string, communityID uuid.UUID) (*community.CommunityResponse, error) {
	com, err := s.repo.GetCommunityByID(ctx, communityID)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("community_id", communityID.String()).Msg("failed to fetch community")
		return nil, err
	}
	return s.withViewerRole(ctx, viewerID, com)
}

func (s *CommunityService) GetBySlug(ctx context.Context, viewerID *string, slug string) (*community.CommunityResponse, error) {
	com, err := s.repo.GetCommunityBySlug(ctx, slug)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("slug", slug).Msg("failed to fetch community")
		return nil, err
	}
	return s.withViewerRole(ctx, viewerID, com)
}

func (s *CommunityService) withViewerRole(ctx context.Context, viewerID *string, com *community.Community) (*community.CommunityResponse, error) {
	resp := &community.CommunityResponse{Community: *com}
	if viewerID == nil {
		return resp, nil
	}

	role, err := s.repo.GetUserRole(ctx, com.ID, *viewerID)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("community_id", com.ID.String()).Str("user_id", *viewerID).Msg("failed to fetch viewer role")
		return nil, err
	}
	resp.ViewerRole = role
	return resp, nil
}

func (s *CommunityService) UpdateSettings(ctx context.Context, userID string, communityID uuid.UUID, payload *community.UpdateCommunitySettingsPayload) (*community.Community, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	payload.UserID = userID

	updated, err := s.repo.UpdateCommunitySettings(ctx, communityID, payload)
	if err != nil {
		logger.Error().Err(err).Str("community_id", communityID.String()).Str("user_id", userID).Msg("failed to update community settings")
		return nil, err
	}

	logger.Info().Str("event", "community_settings_updated").Str("community_id", communityID.String()).Str("user_id", userID).Msg("community settings updated")
	return updated, nil
}

func (s *CommunityService) DeleteCommunity(ctx context.Context, userID string, communityID uuid.UUID) error {
	logger := middleware.GetLoggerFromContext(ctx)

	if err := s.repo.DeleteCommunity(ctx, communityID, userID); err != nil {
		logger.Error().Err(err).Str("community_id", communityID.String()).Str("user_id", userID).Msg("failed to delete community")
		return err
	}

	logger.Info().Str("event", "community_deleted").Str("community_id", communityID.String()).Str("user_id", userID).Msg("community deleted")
	return nil
}

func (s *CommunityService) JoinCommunity(ctx context.Context, userID string, communityID uuid.UUID) (*community.CommunityMember, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	member, err := s.repo.JoinCommunity(ctx, userID, communityID)
	if err != nil {
		logger.Error().Err(err).Str("community_id", communityID.String()).Str("user_id", userID).Msg("failed to join community")
		return nil, err
	}

	logger.Info().Str("event", "community_joined").Str("community_id", communityID.String()).Str("user_id", userID).Msg("user joined community")
	return member, nil
}

func (s *CommunityService) LeaveCommunity(ctx context.Context, userID string, communityID uuid.UUID) error {
	logger := middleware.GetLoggerFromContext(ctx)

	if err := s.repo.LeaveCommunity(ctx, userID, communityID); err != nil {
		logger.Error().Err(err).Str("community_id", communityID.String()).Str("user_id", userID).Msg("failed to leave community")
		return err
	}

	logger.Info().Str("event", "community_left").Str("community_id", communityID.String()).Str("user_id", userID).Msg("user left community")
	return nil
}

func (s *CommunityService) GetMembers(ctx context.Context, viewerID string, communityID uuid.UUID, query *community.GetCommunityMembersQuery) (*model.CursorPaginatedResponse[community.MiniCommunityUser], error) {
	res, err := s.repo.GetCommunityMembers(ctx, viewerID, communityID, query)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("community_id", communityID.String()).Msg("failed to list community members")
		return nil, err
	}
	return res, nil
}

func (s *CommunityService) GetCommunityPostByID(ctx context.Context, viewerID string, communityID uuid.UUID, postID uuid.UUID) (*post.PopulatedPost, error) {
	p, err := s.repo.GetCommunityPostByID(ctx, viewerID, postID, communityID)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("community_id", communityID.String()).Str("post_id", postID.String()).Msg("failed to fetch community post")
		return nil, err
	}
	return p, nil
}

func (s *CommunityService) DeleteCommunityPost(ctx context.Context, userID string, communityID uuid.UUID, postID uuid.UUID) error {
	logger := middleware.GetLoggerFromContext(ctx)

	if err := s.repo.DeleteCommunityPost(ctx, userID, communityID, &community.DeleteCommunityPostPayload{PostID: postID}); err != nil {
		logger.Error().Err(err).Str("community_id", communityID.String()).Str("post_id", postID.String()).Str("moderator_id", userID).Msg("failed to delete community post")
		return err
	}

	logger.Info().Str("event", "community_post_deleted").Str("community_id", communityID.String()).Str("post_id", postID.String()).Str("moderator_id", userID).Msg("community post deleted by moderator")
	return nil
}

func (s *CommunityService) ChangeMemberRole(ctx context.Context, userID string, communityID uuid.UUID, payload *community.ChangeMemberRoleInCommunityPayload) (*community.CommunityMember, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	member, err := s.repo.ChangeMemberRoleInCommunity(ctx, communityID, userID, payload)
	if err != nil {
		logger.Error().Err(err).Str("community_id", communityID.String()).Str("target_user_id", payload.TargetUserID).Msg("failed to change member role")
		return nil, err
	}

	logger.Info().Str("event", "community_member_role_changed").Str("community_id", communityID.String()).Str("target_user_id", payload.TargetUserID).Str("new_role", string(payload.NewRole)).Msg("member role changed")
	return member, nil
}

func (s *CommunityService) BanMember(ctx context.Context, userID string, communityID uuid.UUID, payload *community.BanCommunityMemberPayload) (*community.BannedFromCommunityUser, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	banned, err := s.repo.BanUserFromCommunity(ctx, userID, communityID, payload)
	if err != nil {
		logger.Error().Err(err).Str("community_id", communityID.String()).Str("target_user_id", payload.UserIDToBan).Msg("failed to ban member")
		return nil, err
	}

	logger.Info().Str("event", "community_member_banned").Str("community_id", communityID.String()).Str("target_user_id", payload.UserIDToBan).Str("moderator_id", userID).Msg("member banned")
	return banned, nil
}

func (s *CommunityService) ReportPost(ctx context.Context, userID string, payload *community.ReportCommunityPostPayload) (*community.CommunityReport, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	report, err := s.repo.ReportCommunityPost(ctx, userID, payload)
	if err != nil {
		logger.Error().Err(err).Str("community_id", payload.CommunityID.String()).Str("post_id", payload.PostID.String()).Msg("failed to report post")
		return nil, err
	}

	logger.Info().Str("event", "community_post_reported").Str("report_id", report.ID.String()).Str("post_id", payload.PostID.String()).Str("reporter_id", userID).Msg("post reported")
	return report, nil
}

func (s *CommunityService) GetReports(ctx context.Context, userID string, communityID uuid.UUID, query *community.GetCommunityReportsQuery) (*model.CursorPaginatedResponse[community.CommunityReport], error) {
	res, err := s.repo.GetCommunityReports(ctx, userID, communityID, query)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("community_id", communityID.String()).Msg("failed to list community reports")
		return nil, err
	}
	return res, nil
}

func (s *CommunityService) GetReportByID(ctx context.Context, userID string, communityID uuid.UUID, reportID uuid.UUID) (*community.CommunityReport, error) {
	report, err := s.repo.GetReportByID(ctx, userID, communityID, reportID)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("community_id", communityID.String()).Str("report_id", reportID.String()).Msg("failed to fetch report")
		return nil, err
	}
	return report, nil
}

func (s *CommunityService) ResolveReport(ctx context.Context, userID string, communityID uuid.UUID, payload *community.ResolveCommunityPostReportPayload) (*community.CommunityReport, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	resolved, err := s.repo.ResolveCommunityPostReport(ctx, userID, communityID, payload)
	if err != nil {
		logger.Error().Err(err).Str("community_id", communityID.String()).Str("report_id", payload.ReportID.String()).Msg("failed to resolve report")
		return nil, err
	}

	logger.Info().Str("event", "community_report_resolved").Str("report_id", payload.ReportID.String()).Str("status", string(payload.UpdatedStatus)).Str("moderator_id", userID).Msg("report resolved")
	return resolved, nil
}
