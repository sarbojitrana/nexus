package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/sarbojitrana/nexus/internal/errs"
	"github.com/sarbojitrana/nexus/internal/model/chat"
	"github.com/sarbojitrana/nexus/internal/server"
)

type ChatRepository struct {
	server *server.Server
}

func NewChatRepository(server *server.Server) *ChatRepository {
	return &ChatRepository{server: server}
}

func (r *ChatRepository) GetOrCreateDirectConversation(ctx context.Context, userA, userB string) (*chat.Conversation, error) {
	rows, err := r.server.DB.Pool.Query(ctx, `
		SELECT c.* FROM conversations c
		JOIN conversation_participants cp1 ON cp1.conversation_id = c.id AND cp1.user_id = @user_a
		JOIN conversation_participants cp2 ON cp2.conversation_id = c.id AND cp2.user_id = @user_b
		WHERE c.is_group = false
	`, pgx.NamedArgs{"user_a": userA, "user_b": userB})
	if err != nil {
		return nil, fmt.Errorf("failed to look up direct conversation: %w", err)
	}

	existing, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[chat.Conversation])
	if err == nil {
		return &existing, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to scan direct conversation: %w", err)
	}

	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	crows, err := tx.Query(ctx, `
		INSERT INTO conversations (is_group, created_by) VALUES (false, @created_by) RETURNING *
	`, pgx.NamedArgs{"created_by": userA})
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}
	conv, err := pgx.CollectOneRow(crows, pgx.RowToStructByName[chat.Conversation])
	if err != nil {
		return nil, fmt.Errorf("failed to parse created conversation: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO conversation_participants (conversation_id, user_id) VALUES (@conversation_id, @user_a), (@conversation_id, @user_b)
	`, pgx.NamedArgs{"conversation_id": conv.ID, "user_a": userA, "user_b": userB})
	if err != nil {
		return nil, fmt.Errorf("failed to add participants: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &conv, nil
}

func (r *ChatRepository) CreateGroupConversation(ctx context.Context, creatorID, name string, inviteeIDs []string) (*chat.Conversation, error) {
	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		INSERT INTO conversations (is_group, name, created_by) VALUES (true, @name, @created_by) RETURNING *
	`, pgx.NamedArgs{"name": name, "created_by": creatorID})
	if err != nil {
		return nil, fmt.Errorf("failed to create group conversation: %w", err)
	}
	conv, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[chat.Conversation])
	if err != nil {
		return nil, fmt.Errorf("failed to parse created conversation: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO conversation_participants (conversation_id, user_id) VALUES (@conversation_id, @user_id)
	`, pgx.NamedArgs{"conversation_id": conv.ID, "user_id": creatorID})
	if err != nil {
		return nil, fmt.Errorf("failed to add creator as participant: %w", err)
	}

	for _, inviteeID := range inviteeIDs {
		_, err = tx.Exec(ctx, `
			INSERT INTO conversation_invitations (conversation_id, inviter_id, invitee_id) VALUES (@conversation_id, @inviter_id, @invitee_id)
		`, pgx.NamedArgs{"conversation_id": conv.ID, "inviter_id": creatorID, "invitee_id": inviteeID})
		if err != nil {
			return nil, fmt.Errorf("failed to invite user %s: %w", inviteeID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &conv, nil
}

func (r *ChatRepository) InviteToConversation(ctx context.Context, conversationID uuid.UUID, inviterID, inviteeID string) (*chat.Invitation, error) {
	rows, err := r.server.DB.Pool.Query(ctx, `
		INSERT INTO conversation_invitations (conversation_id, inviter_id, invitee_id)
		VALUES (@conversation_id, @inviter_id, @invitee_id)
		RETURNING *
	`, pgx.NamedArgs{"conversation_id": conversationID, "inviter_id": inviterID, "invitee_id": inviteeID})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			code := "ALREADY_INVITED"
			return nil, errs.NewBadRequestError("user has already been invited", false, &code, nil, nil)
		}
		return nil, fmt.Errorf("failed to invite user %s: %w", inviteeID, err)
	}

	inv, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[chat.Invitation])
	if err != nil {
		return nil, fmt.Errorf("failed to parse invitation: %w", err)
	}
	return &inv, nil
}

func (r *ChatRepository) RespondToInvitation(ctx context.Context, invitationID uuid.UUID, userID string, accept bool) (uuid.UUID, error) {
	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	status := chat.InvitationRejected
	if accept {
		status = chat.InvitationAccepted
	}

	var conversationID uuid.UUID
	err = tx.QueryRow(ctx, `
		UPDATE conversation_invitations
		SET status = @status, responded_at = CURRENT_TIMESTAMP
		WHERE id = @id AND invitee_id = @user_id AND status = 'pending'
		RETURNING conversation_id
	`, pgx.NamedArgs{"status": status, "id": invitationID, "user_id": userID}).Scan(&conversationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			code := "INVITATION_NOT_FOUND"
			return uuid.Nil, errs.NewNotFoundError("invitation not found or already responded to", false, &code)
		}
		return uuid.Nil, fmt.Errorf("failed to update invitation: %w", err)
	}

	if accept {
		_, err = tx.Exec(ctx, `
			INSERT INTO conversation_participants (conversation_id, user_id)
			VALUES (@conversation_id, @user_id)
			ON CONFLICT DO NOTHING
		`, pgx.NamedArgs{"conversation_id": conversationID, "user_id": userID})
		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to add participant: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return conversationID, nil
}

func (r *ChatRepository) GetPendingInvitations(ctx context.Context, userID string) ([]chat.Invitation, error) {
	rows, err := r.server.DB.Pool.Query(ctx, `
		SELECT * FROM conversation_invitations WHERE invitee_id = @user_id AND status = 'pending' ORDER BY created_at DESC
	`, pgx.NamedArgs{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to query invitations: %w", err)
	}

	invitations, err := pgx.CollectRows(rows, pgx.RowToStructByName[chat.Invitation])
	if err != nil {
		return nil, fmt.Errorf("failed to parse invitations: %w", err)
	}
	return invitations, nil
}

func (r *ChatRepository) IsParticipant(ctx context.Context, conversationID uuid.UUID, userID string) (bool, error) {
	var exists bool
	err := r.server.DB.Pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM conversation_participants WHERE conversation_id = @conversation_id AND user_id = @user_id)
	`, pgx.NamedArgs{"conversation_id": conversationID, "user_id": userID}).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check participant status: %w", err)
	}
	return exists, nil
}

func (r *ChatRepository) ListConversationsForUser(ctx context.Context, userID string) ([]chat.Conversation, error) {
	rows, err := r.server.DB.Pool.Query(ctx, `
		SELECT c.* FROM conversations c
		JOIN conversation_participants cp ON cp.conversation_id = c.id
		WHERE cp.user_id = @user_id
		ORDER BY c.updated_at DESC
		LIMIT 50
	`, pgx.NamedArgs{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}

	conversations, err := pgx.CollectRows(rows, pgx.RowToStructByName[chat.Conversation])
	if err != nil {
		return nil, fmt.Errorf("failed to parse conversations: %w", err)
	}
	return conversations, nil
}

func (r *ChatRepository) GetParticipants(ctx context.Context, conversationID uuid.UUID) ([]chat.Participant, error) {
	rows, err := r.server.DB.Pool.Query(ctx, `
		SELECT * FROM conversation_participants WHERE conversation_id = @conversation_id
	`, pgx.NamedArgs{"conversation_id": conversationID})
	if err != nil {
		return nil, fmt.Errorf("failed to query participants: %w", err)
	}

	participants, err := pgx.CollectRows(rows, pgx.RowToStructByName[chat.Participant])
	if err != nil {
		return nil, fmt.Errorf("failed to parse participants: %w", err)
	}
	return participants, nil
}

func (r *ChatRepository) GetLastMessage(ctx context.Context, conversationID uuid.UUID) (*chat.Message, error) {
	rows, err := r.server.DB.Pool.Query(ctx, `
		SELECT * FROM messages WHERE conversation_id = @conversation_id ORDER BY created_at DESC LIMIT 1
	`, pgx.NamedArgs{"conversation_id": conversationID})
	if err != nil {
		return nil, fmt.Errorf("failed to query last message: %w", err)
	}

	msg, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[chat.Message])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to parse last message: %w", err)
	}
	return &msg, nil
}

func (r *ChatRepository) GetUnreadCount(ctx context.Context, conversationID uuid.UUID, userID string) (int, error) {
	var count int
	err := r.server.DB.Pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM messages m
		JOIN conversation_participants cp ON cp.conversation_id = m.conversation_id AND cp.user_id = @user_id
		WHERE m.conversation_id = @conversation_id
			AND m.sender_id != @user_id
			AND (cp.last_read_at IS NULL OR m.created_at > cp.last_read_at)
	`, pgx.NamedArgs{"conversation_id": conversationID, "user_id": userID}).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count unread messages: %w", err)
	}
	return count, nil
}

func (r *ChatRepository) GetMessages(ctx context.Context, conversationID uuid.UUID, query *chat.GetMessagesQuery) ([]chat.Message, bool, error) {
	stmt := `SELECT * FROM messages WHERE conversation_id = @conversation_id`
	args := pgx.NamedArgs{"conversation_id": conversationID}

	if query.CursorCreatedAt != nil {
		args["cursor_created_at"] = *query.CursorCreatedAt
		stmt += " AND created_at < @cursor_created_at"
	}

	limit := 30
	args["limit_plus_one"] = limit + 1
	stmt += " ORDER BY created_at DESC LIMIT @limit_plus_one"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, false, fmt.Errorf("failed to query messages: %w", err)
	}

	messages, err := pgx.CollectRows(rows, pgx.RowToStructByName[chat.Message])
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse messages: %w", err)
	}

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}
	return messages, hasMore, nil
}

func (r *ChatRepository) SendMessage(ctx context.Context, conversationID uuid.UUID, senderID string, payload *chat.SendMessagePayload) (*chat.Message, error) {
	rows, err := r.server.DB.Pool.Query(ctx, `
		INSERT INTO messages (conversation_id, sender_id, content, attachment_key, attachment_mime_type, attachment_file_size)
		VALUES (@conversation_id, @sender_id, @content, @attachment_key, @attachment_mime_type, @attachment_file_size)
		RETURNING *
	`, pgx.NamedArgs{
		"conversation_id":      conversationID,
		"sender_id":            senderID,
		"content":              payload.Content,
		"attachment_key":       payload.AttachmentKey,
		"attachment_mime_type": payload.AttachmentMimeType,
		"attachment_file_size": payload.AttachmentFileSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	msg, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[chat.Message])
	if err != nil {
		return nil, fmt.Errorf("failed to parse sent message: %w", err)
	}
	return &msg, nil
}

func (r *ChatRepository) MarkRead(ctx context.Context, conversationID uuid.UUID, userID string) error {
	result, err := r.server.DB.Pool.Exec(ctx, `
		UPDATE conversation_participants SET last_read_at = CURRENT_TIMESTAMP
		WHERE conversation_id = @conversation_id AND user_id = @user_id
	`, pgx.NamedArgs{"conversation_id": conversationID, "user_id": userID})
	if err != nil {
		return fmt.Errorf("failed to mark conversation read: %w", err)
	}
	if result.RowsAffected() == 0 {
		code := "NOT_A_PARTICIPANT"
		return errs.NewNotFoundError("not a participant of this conversation", false, &code)
	}
	return nil
}
