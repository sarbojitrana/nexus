package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/follow"
	"github.com/sarbojitrana/nexus/internal/model/notification"
	"github.com/sarbojitrana/nexus/internal/model/user"
	"github.com/sarbojitrana/nexus/internal/repository"
	"github.com/sarbojitrana/nexus/internal/server"
)

type FollowService struct {
	server        *server.Server
	repo          *repository.FollowRepository
	notifications *NotificationService
}

func NewFollowService(s *server.Server, repo *repository.FollowRepository, notifications *NotificationService) *FollowService {
	return &FollowService{
		server:        s,
		repo:          repo,
		notifications: notifications,
	}
}

func (s *FollowService) FollowCommunity(ctx context.Context, userID string, payload *follow.FollowCommunityPayload) (*follow.CommunityFollow, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	f, err := s.repo.FollowCommunity(ctx, userID, payload)
	if err != nil {
		logger.Error().Err(err).Str("community_id", payload.CommunityID.String()).Str("user_id", userID).Msg("failed to follow community")
		return nil, err
	}

	logger.Info().Str("event", "community_followed").Str("community_id", payload.CommunityID.String()).Str("user_id", userID).Msg("user followed community")
	return f, nil
}

func (s *FollowService) UnFollowCommunity(ctx context.Context, userID string, payload *follow.UnFollowCommunityPayload) error {
	logger := middleware.GetLoggerFromContext(ctx)

	if err := s.repo.UnFollowCommunity(ctx, userID, payload); err != nil {
		logger.Error().Err(err).Str("community_id", payload.CommunityID.String()).Str("user_id", userID).Msg("failed to unfollow community")
		return err
	}

	logger.Info().Str("event", "community_unfollowed").Str("community_id", payload.CommunityID.String()).Str("user_id", userID).Msg("user unfollowed community")
	return nil
}

func (s *FollowService) IsFollowingCommunity(ctx context.Context, userID string, communityID uuid.UUID) (*bool, error) {
	res, err := s.repo.IsFollowingCommunity(ctx, userID, communityID)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("community_id", communityID.String()).Str("user_id", userID).Msg("failed to check community follow status")
		return nil, err
	}
	return res, nil
}

func (s *FollowService) FollowUser(ctx context.Context, userID string, payload *follow.FollowUserPayload) (*follow.UserFollow, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	f, err := s.repo.FollowUser(ctx, userID, payload)
	if err != nil {
		logger.Error().Err(err).Str("following_id", payload.FollowingID).Str("user_id", userID).Msg("failed to follow user")
		return nil, err
	}

	logger.Info().Str("event", "user_followed").Str("following_id", payload.FollowingID).Str("user_id", userID).Msg("user followed another user")
	s.notifications.Notify(ctx, payload.FollowingID, &userID, notification.TypeFollow, map[string]string{"followerId": userID})
	return f, nil
}

func (s *FollowService) UnFollowUser(ctx context.Context, userID string, payload *follow.UnFollowUserPayload) error {
	logger := middleware.GetLoggerFromContext(ctx)

	if err := s.repo.UnFollowUser(ctx, userID, payload); err != nil {
		logger.Error().Err(err).Str("following_id", payload.FollowingID).Str("user_id", userID).Msg("failed to unfollow user")
		return err
	}

	logger.Info().Str("event", "user_unfollowed").Str("following_id", payload.FollowingID).Str("user_id", userID).Msg("user unfollowed another user")
	return nil
}

func (s *FollowService) IsFollowingUser(ctx context.Context, userID string, targetUserID string) (*bool, error) {
	res, err := s.repo.IsFollowingUser(ctx, userID, targetUserID)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("target_user_id", targetUserID).Str("user_id", userID).Msg("failed to check user follow status")
		return nil, err
	}
	return res, nil
}

func (s *FollowService) GetFollowers(ctx context.Context, viewerID *string, userID string, query *follow.GetFollowersQuery) (*model.CursorPaginatedResponse[user.MiniUser], error) {
	res, err := s.repo.GetFollowers(ctx, viewerID, userID, query)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("user_id", userID).Msg("failed to list followers")
		return nil, err
	}
	return res, nil
}

func (s *FollowService) GetFollowing(ctx context.Context, viewerID *string, userID string, query *follow.GetFollowersQuery) (*model.CursorPaginatedResponse[user.MiniUser], error) {
	res, err := s.repo.GetFollowing(ctx, viewerID, userID, query)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("user_id", userID).Msg("failed to list following")
		return nil, err
	}
	return res, nil
}
