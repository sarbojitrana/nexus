package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sarbojitrana/nexus/internal/handler"
	"github.com/sarbojitrana/nexus/internal/middleware"
)

func registerChatRoutes(g *echo.Group, h *handler.Handlers, mw *middleware.Middlewares) {
	g.POST("/conversations/direct", h.Chat.StartDirectConversation, mw.Auth.RequireAuth)
	g.POST("/conversations/group", h.Chat.CreateGroup, mw.Auth.RequireAuth)
	g.GET("/conversations", h.Chat.ListConversations, mw.Auth.RequireAuth)
	g.POST("/conversations/:id/invitations", h.Chat.InviteToConversation, mw.Auth.RequireAuth)
	g.GET("/conversations/:id/messages", h.Chat.GetMessages, mw.Auth.RequireAuth)
	g.POST("/conversations/:id/messages", h.Chat.SendMessage, mw.Auth.RequireAuth)
	g.POST("/conversations/:id/read", h.Chat.MarkRead, mw.Auth.RequireAuth)

	g.GET("/invitations", h.Chat.GetPendingInvitations, mw.Auth.RequireAuth)
	g.POST("/invitations/:id/accept", h.Chat.AcceptInvitation, mw.Auth.RequireAuth)
	g.POST("/invitations/:id/reject", h.Chat.RejectInvitation, mw.Auth.RequireAuth)
}

func registerChatWebSocketRoute(r *echo.Echo, h *handler.Handlers) {
	r.GET("/ws/chat", h.Chat.ServeWebSocket)
}
