package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/notification"
	"github.com/sarbojitrana/nexus/internal/server"
)

type NotificationRepository struct {
	server *server.Server
}

func NewNotificationRepository(server *server.Server) *NotificationRepository {
	return &NotificationRepository{server: server}
}

func (r *NotificationRepository) Create(ctx context.Context, userID string, actorID *string, notifType notification.Type, data any) (*notification.Notification, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal notification data: %w", err)
	}

	stmt := `
		INSERT INTO notifications (user_id, actor_id, type, data)
		VALUES (@user_id, @actor_id, @type, @data)
		RETURNING *
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id": userID, "actor_id": actorID, "type": notifType, "data": dataJSON,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}
	defer rows.Close()

	n, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[notification.Notification])
	if err != nil {
		return nil, fmt.Errorf("failed to parse notification: %w", err)
	}
	return &n, nil
}

func (r *NotificationRepository) List(ctx context.Context, userID string, query *notification.GetNotificationsQuery) (*model.CursorPaginatedResponse[notification.Notification], error) {
	stmt := `SELECT * FROM notifications WHERE user_id = @user_id`
	args := pgx.NamedArgs{"user_id": userID}

	if query.CursorCreatedAt != nil {
		args["cursor_created_at"] = *query.CursorCreatedAt
		stmt += " AND created_at <= @cursor_created_at"
	}

	limit := 20
	args["limit_plus_one"] = limit + 1
	stmt += " ORDER BY created_at DESC LIMIT @limit_plus_one"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}

	notifications, err := pgx.CollectRows(rows, pgx.RowToStructByName[notification.Notification])
	if err != nil {
		return nil, fmt.Errorf("failed to parse notifications: %w", err)
	}

	if len(notifications) < limit+1 {
		var cursorCreatedAt time.Time
		if query.CursorCreatedAt != nil {
			cursorCreatedAt = *query.CursorCreatedAt
		}
		return &model.CursorPaginatedResponse[notification.Notification]{
			Data:            notifications,
			CursorCreatedAt: cursorCreatedAt,
			HasMore:         false,
		}, nil
	}

	return &model.CursorPaginatedResponse[notification.Notification]{
		Data:            notifications[:limit],
		CursorCreatedAt: notifications[limit].CreatedAt,
		HasMore:         true,
	}, nil
}

func (r *NotificationRepository) MarkRead(ctx context.Context, id uuid.UUID, userID string) error {
	_, err := r.server.DB.Pool.Exec(ctx, `
		UPDATE notifications SET is_read = true WHERE id = @id AND user_id = @user_id
	`, pgx.NamedArgs{"id": id, "user_id": userID})
	if err != nil {
		return fmt.Errorf("failed to mark notification read: %w", err)
	}
	return nil
}

func (r *NotificationRepository) MarkAllRead(ctx context.Context, userID string) error {
	_, err := r.server.DB.Pool.Exec(ctx, `
		UPDATE notifications SET is_read = true WHERE user_id = @user_id AND is_read = false
	`, pgx.NamedArgs{"user_id": userID})
	if err != nil {
		return fmt.Errorf("failed to mark notifications read: %w", err)
	}
	return nil
}

func (r *NotificationRepository) MarkMessageNotificationsRead(ctx context.Context, userID string, conversationID uuid.UUID) error {
	_, err := r.server.DB.Pool.Exec(ctx, `
		UPDATE notifications
		SET is_read = true
		WHERE user_id = @user_id
			AND type = 'message'
			AND is_read = false
			AND (data ->> 'conversationId')::uuid = @conversation_id
	`, pgx.NamedArgs{"user_id": userID, "conversation_id": conversationID})
	if err != nil {
		return fmt.Errorf("failed to mark message notifications read: %w", err)
	}
	return nil
}
