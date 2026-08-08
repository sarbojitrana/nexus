package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sarbojitrana/nexus/internal/handler"
	"github.com/sarbojitrana/nexus/internal/middleware"
)

func registerUserRoutes(g *echo.Group, h *handler.Handlers, mw *middleware.Middlewares) {
	g.GET("/users", h.User.GetUsers, mw.Auth.OptionalAuth)
	g.GET("/users/:id", h.User.GetUserByID, mw.Auth.OptionalAuth)
	g.GET("/users/:id/posts", h.User.GetPostsByUserID, mw.Auth.OptionalAuth)

	g.GET("/users/me", h.User.GetMe, mw.Auth.RequireAuth)
	g.PATCH("/users/me", h.User.UpdateMe, mw.Auth.RequireAuth)
	g.DELETE("/users/me", h.User.DeleteMe, mw.Auth.RequireAuth)

	g.GET("/users/me/settings", h.User.GetMySettings, mw.Auth.RequireAuth)
	g.PATCH("/users/me/settings", h.User.UpdateMySettings, mw.Auth.RequireAuth)

	g.POST("/users/me/password-changed", h.User.NotifyPasswordChanged, mw.Auth.RequireAuth)

	g.GET("/users/me/blocks", h.User.GetBlockedUsers, mw.Auth.RequireAuth)
	g.POST("/users/:id/block", h.User.BlockUser, mw.Auth.RequireAuth)
	g.DELETE("/users/:id/block", h.User.UnblockUser, mw.Auth.RequireAuth)
	g.GET("/users/:id/block", h.User.IsBlockingUser, mw.Auth.RequireAuth)
}
