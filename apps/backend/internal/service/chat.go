package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/sarbojitrana/nexus/internal/errs"
	"github.com/sarbojitrana/nexus/internal/lib/ws"
	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/model/chat"
	"github.com/sarbojitrana/nexus/internal/model/notification"
	"github.com/sarbojitrana/nexus/internal/model/user"
	"github.com/sarbojitrana/nexus/internal/repository"
	"github.com/sarbojitrana/nexus/internal/server"
)

type ChatService struct {
	server        *server.Server
	repo          *repository.ChatRepository
	userRepo      *repository.UserRepository
	followRepo    *repository.FollowRepository
	hub           *ws.Hub
	notifications *NotificationService
}

func NewChatService(
	s *server.Server,
	repo *repository.ChatRepository,
	userRepo *repository.UserRepository,
	followRepo *repository.FollowRepository,
	hub *ws.Hub,
	notifications *NotificationService,
) *ChatService {
	return &ChatService{
		server:        s,
		repo:          repo,
		userRepo:      userRepo,
		followRepo:    followRepo,
		hub:           hub,
		notifications: notifications,
	}
}

func (s *ChatService) StartDirectConversation(ctx context.Context, userID, targetUserID string) (*chat.ConversationSummary, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	if userID == targetUserID {
		code := "CANNOT_MESSAGE_SELF"
		return nil, errs.NewBadRequestError("cannot start a conversation with yourself", false, &code, nil, nil)
	}

	conv, err := s.repo.GetOrCreateDirectConversation(ctx, userID, targetUserID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Str("target_user_id", targetUserID).Msg("failed to start direct conversation")
		return nil, err
	}

	return s.buildSummary(ctx, userID, conv)
}

func (s *ChatService) CreateGroup(ctx context.Context, creatorID string, payload *chat.CreateGroupConversationPayload) (*chat.ConversationSummary, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	allowedInvitees := make([]string, 0, len(payload.InviteeIDs))
	for _, inviteeID := range payload.InviteeIDs {
		if inviteeID == creatorID {
			continue
		}
		allowed, err := s.canInvite(ctx, creatorID, inviteeID)
		if err != nil {
			logger.Error().Err(err).Str("invitee_id", inviteeID).Msg("failed to check invite permission, skipping invitee")
			continue
		}
		if allowed {
			allowedInvitees = append(allowedInvitees, inviteeID)
		}
	}

	if len(allowedInvitees) == 0 {
		code := "NO_INVITABLE_USERS"
		return nil, errs.NewBadRequestError("none of the selected users can be invited", false, &code, nil, nil)
	}

	conv, err := s.repo.CreateGroupConversation(ctx, creatorID, payload.Name, allowedInvitees)
	if err != nil {
		logger.Error().Err(err).Str("user_id", creatorID).Msg("failed to create group conversation")
		return nil, err
	}

	logger.Info().Str("event", "group_conversation_created").Str("conversation_id", conv.ID.String()).Str("user_id", creatorID).Msg("group conversation created")

	for _, inviteeID := range allowedInvitees {
		s.notifications.Notify(ctx, inviteeID, &creatorID, notification.TypeGroupInvitation, map[string]string{"conversationId": conv.ID.String()})
	}

	return s.buildSummary(ctx, creatorID, conv)
}

func (s *ChatService) canInvite(ctx context.Context, inviterID, inviteeID string) (bool, error) {
	invitee, err := s.userRepo.GetUserByID(ctx, nil, inviteeID)
	if err != nil {
		return false, err
	}

	switch invitee.GroupInvitePermission {
	case user.GroupInvitePermissionNoOne:
		return false, nil
	case user.GroupInvitePermissionFollowersOnly:
		isFollowing, err := s.followRepo.IsFollowingUser(ctx, inviteeID, inviterID)
		if err != nil {
			return false, err
		}
		return *isFollowing, nil
	default: // everyone
		return true, nil
	}
}

func (s *ChatService) InviteToConversation(ctx context.Context, conversationID uuid.UUID, inviterID string, payload *chat.InviteToConversationPayload) (*chat.Invitation, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	isParticipant, err := s.repo.IsParticipant(ctx, conversationID, inviterID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, errs.NewForbiddenError("must be a participant to invite others", false)
	}

	allowed, err := s.canInvite(ctx, inviterID, payload.UserID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to check invite permission")
		return nil, err
	}
	if !allowed {
		code := "CANNOT_INVITE_USER"
		return nil, errs.NewBadRequestError("this user can't be invited", false, &code, nil, nil)
	}

	inv, err := s.repo.InviteToConversation(ctx, conversationID, inviterID, payload.UserID)
	if err != nil {
		logger.Error().Err(err).Str("conversation_id", conversationID.String()).Str("invitee_id", payload.UserID).Msg("failed to invite user to conversation")
		return nil, err
	}

	logger.Info().Str("event", "conversation_invitation_sent").Str("conversation_id", conversationID.String()).Str("invitee_id", payload.UserID).Msg("group invitation sent")
	s.notifications.Notify(ctx, payload.UserID, &inviterID, notification.TypeGroupInvitation, map[string]string{"conversationId": conversationID.String()})

	return inv, nil
}

func (s *ChatService) RespondToInvitation(ctx context.Context, invitationID uuid.UUID, userID string, accept bool) error {
	logger := middleware.GetLoggerFromContext(ctx)

	if _, err := s.repo.RespondToInvitation(ctx, invitationID, userID, accept); err != nil {
		logger.Error().Err(err).Str("invitation_id", invitationID.String()).Msg("failed to respond to invitation")
		return err
	}

	event := "invitation_rejected"
	if accept {
		event = "invitation_accepted"
	}
	logger.Info().Str("event", event).Str("invitation_id", invitationID.String()).Str("user_id", userID).Msg("invitation responded to")
	return nil
}

func (s *ChatService) GetPendingInvitations(ctx context.Context, userID string) ([]chat.Invitation, error) {
	invitations, err := s.repo.GetPendingInvitations(ctx, userID)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("user_id", userID).Msg("failed to fetch pending invitations")
		return nil, err
	}
	return invitations, nil
}

func (s *ChatService) ListConversations(ctx context.Context, userID string) ([]chat.ConversationSummary, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	conversations, err := s.repo.ListConversationsForUser(ctx, userID)
	if err != nil {
		logger.Error().Err(err).Str("user_id", userID).Msg("failed to list conversations")
		return nil, err
	}

	summaries := make([]chat.ConversationSummary, 0, len(conversations))
	for i := range conversations {
		summary, err := s.buildSummary(ctx, userID, &conversations[i])
		if err != nil {
			logger.Error().Err(err).Str("conversation_id", conversations[i].ID.String()).Msg("failed to build conversation summary, skipping")
			continue
		}
		summaries = append(summaries, *summary)
	}
	return summaries, nil
}

func (s *ChatService) buildSummary(ctx context.Context, viewerID string, conv *chat.Conversation) (*chat.ConversationSummary, error) {
	participants, err := s.repo.GetParticipants(ctx, conv.ID)
	if err != nil {
		return nil, err
	}

	views := make([]chat.ParticipantView, 0, len(participants))
	for _, p := range participants {
		u, err := s.userRepo.GetUserByID(ctx, nil, p.UserID)
		if err != nil {
			continue
		}

		view := chat.ParticipantView{
			UserID:      p.UserID,
			Username:    u.Username,
			DisplayName: u.DisplayName,
			AvatarKey:   u.AvatarKey,
			IsOnline:    u.ShowOnlineStatus && s.hub != nil && s.hub.IsOnline(p.UserID),
		}
		if p.UserID == viewerID || u.ShareReadReceipts {
			view.LastReadAt = p.LastReadAt
		}
		views = append(views, view)
	}

	lastMessage, err := s.repo.GetLastMessage(ctx, conv.ID)
	if err != nil {
		return nil, err
	}

	unread, err := s.repo.GetUnreadCount(ctx, conv.ID, viewerID)
	if err != nil {
		return nil, err
	}

	return &chat.ConversationSummary{
		Conversation: *conv,
		LastMessage:  lastMessage,
		UnreadCount:  unread,
		Participants: views,
	}, nil
}

func (s *ChatService) GetMessages(ctx context.Context, conversationID uuid.UUID, userID string, query *chat.GetMessagesQuery) ([]chat.Message, bool, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	isParticipant, err := s.repo.IsParticipant(ctx, conversationID, userID)
	if err != nil {
		return nil, false, err
	}
	if !isParticipant {
		code := "CONVERSATION_NOT_FOUND"
		return nil, false, errs.NewNotFoundError("conversation not found", false, &code)
	}

	messages, hasMore, err := s.repo.GetMessages(ctx, conversationID, query)
	if err != nil {
		logger.Error().Err(err).Str("conversation_id", conversationID.String()).Msg("failed to fetch messages")
		return nil, false, err
	}

	if err := s.notifications.MarkConversationRead(ctx, userID, conversationID); err != nil {
		logger.Error().Err(err).Msg("failed to resolve message notifications for opened conversation")
	}

	return messages, hasMore, nil
}

func (s *ChatService) SendMessage(ctx context.Context, conversationID uuid.UUID, senderID string, payload *chat.SendMessagePayload) (*chat.Message, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	isParticipant, err := s.repo.IsParticipant(ctx, conversationID, senderID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		code := "CONVERSATION_NOT_FOUND"
		return nil, errs.NewNotFoundError("conversation not found", false, &code)
	}

	msg, err := s.repo.SendMessage(ctx, conversationID, senderID, payload)
	if err != nil {
		logger.Error().Err(err).Str("conversation_id", conversationID.String()).Msg("failed to send message")
		return nil, err
	}

	logger.Info().Str("event", "message_sent").Str("conversation_id", conversationID.String()).Str("message_id", msg.ID.String()).Str("user_id", senderID).Msg("message sent")

	participants, err := s.repo.GetParticipants(ctx, conversationID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch participants for broadcast, message was still sent")
		return msg, nil
	}

	recipientIDs := make([]string, 0, len(participants))
	for _, p := range participants {
		recipientIDs = append(recipientIDs, p.UserID)
		if p.UserID == senderID {
			continue
		}
		s.notifications.Notify(ctx, p.UserID, &senderID, notification.TypeMessage, map[string]string{"conversationId": conversationID.String()})
	}

	if s.hub != nil {
		s.hub.SendToUsers(recipientIDs, ws.Event{Type: "message", Payload: msg})
	}

	return msg, nil
}

func (s *ChatService) MarkRead(ctx context.Context, conversationID uuid.UUID, userID string) error {
	if err := s.repo.MarkRead(ctx, conversationID, userID); err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("conversation_id", conversationID.String()).Msg("failed to mark conversation read")
		return err
	}
	return nil
}

func (s *ChatService) BroadcastPresence(userID string, online bool) {
	ctx := context.Background()

	u, err := s.userRepo.GetUserByID(ctx, nil, userID)
	if err != nil || !u.ShowOnlineStatus {
		return
	}

	conversations, err := s.repo.ListConversationsForUser(ctx, userID)
	if err != nil {
		s.server.Logger.Error().Err(err).Str("user_id", userID).Msg("failed to list conversations for presence broadcast")
		return
	}

	seen := make(map[string]bool)
	event := ws.Event{Type: "presence", Payload: map[string]any{"userId": userID, "online": online}}

	for _, conv := range conversations {
		participants, err := s.repo.GetParticipants(ctx, conv.ID)
		if err != nil {
			continue
		}
		for _, p := range participants {
			if p.UserID == userID || seen[p.UserID] {
				continue
			}
			seen[p.UserID] = true
			s.hub.SendToUser(p.UserID, event)
		}
	}
}
