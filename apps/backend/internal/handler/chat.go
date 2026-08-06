package handler

import (
	"encoding/json"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/sarbojitrana/nexus/internal/lib/ws"
	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/chat"
	"github.com/sarbojitrana/nexus/internal/server"
	"github.com/sarbojitrana/nexus/internal/service"
)

type ChatHandler struct {
	Handler
	services *service.Services
	upgrader websocket.Upgrader
}

func NewChatHandler(s *server.Server, services *service.Services) *ChatHandler {
	return &ChatHandler{
		Handler:  NewHandler(s),
		services: services,
		upgrader: ws.NewUpgrader(s.Config.Server.CORSAllowedOrigins),
	}
}

func (h *ChatHandler) StartDirectConversation(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *chat.StartDirectConversationPayload) (*chat.ConversationSummary, error) {
		return h.services.Chat.StartDirectConversation(c.Request().Context(), middleware.GetUserID(c), req.UserID)
	}, http.StatusOK, &chat.StartDirectConversationPayload{})(c)
}

func (h *ChatHandler) CreateGroup(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *chat.CreateGroupConversationPayload) (*chat.ConversationSummary, error) {
		return h.services.Chat.CreateGroup(c.Request().Context(), middleware.GetUserID(c), req)
	}, http.StatusCreated, &chat.CreateGroupConversationPayload{})(c)
}

func (h *ChatHandler) ListConversations(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *model.Empty) ([]chat.ConversationSummary, error) {
		return h.services.Chat.ListConversations(c.Request().Context(), middleware.GetUserID(c))
	}, http.StatusOK, &model.Empty{})(c)
}

func (h *ChatHandler) InviteToConversation(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *chat.InviteToConversationPayload) (*chat.Invitation, error) {
		conversationID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		return h.services.Chat.InviteToConversation(c.Request().Context(), conversationID, middleware.GetUserID(c), req)
	}, http.StatusCreated, &chat.InviteToConversationPayload{})(c)
}

func (h *ChatHandler) GetPendingInvitations(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *model.Empty) ([]chat.Invitation, error) {
		return h.services.Chat.GetPendingInvitations(c.Request().Context(), middleware.GetUserID(c))
	}, http.StatusOK, &model.Empty{})(c)
}

func (h *ChatHandler) AcceptInvitation(c echo.Context) error {
	return HandleNoContent(h.Handler, func(c echo.Context, req *model.Empty) error {
		invitationID, err := parseUUIDParam(c, "id")
		if err != nil {
			return err
		}
		return h.services.Chat.RespondToInvitation(c.Request().Context(), invitationID, middleware.GetUserID(c), true)
	}, http.StatusNoContent, &model.Empty{})(c)
}

func (h *ChatHandler) RejectInvitation(c echo.Context) error {
	return HandleNoContent(h.Handler, func(c echo.Context, req *model.Empty) error {
		invitationID, err := parseUUIDParam(c, "id")
		if err != nil {
			return err
		}
		return h.services.Chat.RespondToInvitation(c.Request().Context(), invitationID, middleware.GetUserID(c), false)
	}, http.StatusNoContent, &model.Empty{})(c)
}

func (h *ChatHandler) GetMessages(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *chat.GetMessagesQuery) (*model.CursorPaginatedResponse[chat.Message], error) {
		conversationID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		messages, hasMore, err := h.services.Chat.GetMessages(c.Request().Context(), conversationID, middleware.GetUserID(c), req)
		if err != nil {
			return nil, err
		}

		var cursor model.CursorPaginatedResponse[chat.Message]
		cursor.Data = messages
		cursor.HasMore = hasMore
		if len(messages) > 0 {
			cursor.CursorCreatedAt = messages[len(messages)-1].CreatedAt
		}
		return &cursor, nil
	}, http.StatusOK, &chat.GetMessagesQuery{})(c)
}

func (h *ChatHandler) SendMessage(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *chat.SendMessagePayload) (*chat.Message, error) {
		conversationID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		return h.services.Chat.SendMessage(c.Request().Context(), conversationID, middleware.GetUserID(c), req)
	}, http.StatusCreated, &chat.SendMessagePayload{})(c)
}

func (h *ChatHandler) MarkRead(c echo.Context) error {
	return HandleNoContent(h.Handler, func(c echo.Context, req *model.Empty) error {
		conversationID, err := parseUUIDParam(c, "id")
		if err != nil {
			return err
		}
		return h.services.Chat.MarkRead(c.Request().Context(), conversationID, middleware.GetUserID(c))
	}, http.StatusNoContent, &model.Empty{})(c)
}

type wsInboundMessage struct {
	Type           string          `json:"type"`
	ConversationID string          `json:"conversationId"`
	Payload        json.RawMessage `json:"payload"`
}

func (h *ChatHandler) ServeWebSocket(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
	}

	claims, err := jwt.Verify(c.Request().Context(), &jwt.VerifyParams{Token: token})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}
	userID := claims.Subject

	conn, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		h.server.Logger.Error().Err(err).Str("user_id", userID).Msg("failed to upgrade websocket connection")
		return nil
	}

	client := ws.NewClient(h.server.Hub, conn, userID)
	h.server.Hub.Register(client)

	go client.WritePump()

	ctx := c.Request().Context()
	client.ReadPump(func(raw []byte) {
		var msg wsInboundMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			return
		}
		if msg.Type != "message" {
			return
		}

		conversationID, err := uuid.Parse(msg.ConversationID)
		if err != nil {
			return
		}

		var payload chat.SendMessagePayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return
		}
		if err := payload.Validate(); err != nil {
			return
		}

		if _, err := h.services.Chat.SendMessage(ctx, conversationID, userID, &payload); err != nil {
			h.server.Logger.Error().Err(err).Str("user_id", userID).Msg("failed to send message over websocket")
		}
	})

	return nil
}
