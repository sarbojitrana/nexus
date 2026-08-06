package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sarbojitrana/nexus/internal/handler"
	"github.com/sarbojitrana/nexus/internal/middleware"
)

func registerNotificationRoutes(g *echo.Group, h *handler.Handlers, mw *middleware.Middlewares) {
	g.GET("/notifications", h.Notification.List, mw.Auth.RequireAuth)
	g.POST("/notifications/:id/read", h.Notification.MarkRead, mw.Auth.RequireAuth)
	g.POST("/notifications/read-all", h.Notification.MarkAllRead, mw.Auth.RequireAuth)
}
