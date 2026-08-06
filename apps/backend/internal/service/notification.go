package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/sarbojitrana/nexus/internal/lib/ws"
	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/notification"
	"github.com/sarbojitrana/nexus/internal/repository"
	"github.com/sarbojitrana/nexus/internal/server"
)

type NotificationService struct {
	server *server.Server
	repo   *repository.NotificationRepository
	hub    *ws.Hub
}

func NewNotificationService(s *server.Server, repo *repository.NotificationRepository, hub *ws.Hub) *NotificationService {
	return &NotificationService{server: s, repo: repo, hub: hub}
}

func (s *NotificationService) List(ctx context.Context, userID string, query *notification.GetNotificationsQuery) (*model.CursorPaginatedResponse[notification.Notification], error) {
	res, err := s.repo.List(ctx, userID, query)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("user_id", userID).Msg("failed to list notifications")
		return nil, err
	}
	return res, nil
}

func (s *NotificationService) MarkRead(ctx context.Context, id uuid.UUID, userID string) error {
	if err := s.repo.MarkRead(ctx, id, userID); err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("notification_id", id.String()).Msg("failed to mark notification read")
		return err
	}
	return nil
}

func (s *NotificationService) MarkAllRead(ctx context.Context, userID string) error {
	if err := s.repo.MarkAllRead(ctx, userID); err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("user_id", userID).Msg("failed to mark all notifications read")
		return err
	}
	return nil
}

func (s *NotificationService) MarkConversationRead(ctx context.Context, userID string, conversationID uuid.UUID) error {
	return s.repo.MarkMessageNotificationsRead(ctx, userID, conversationID)
}

func (s *NotificationService) Notify(ctx context.Context, userID string, actorID *string, notifType notification.Type, data any) {
	logger := middleware.GetLoggerFromContext(ctx)

	n, err := s.repo.Create(ctx, userID, actorID, notifType, data)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Str("type", string(notifType)).Msg("failed to create notification")
		return
	}

	if s.hub != nil {
		s.hub.SendToUser(userID, ws.Event{Type: "notification", Payload: n})
	}
}
