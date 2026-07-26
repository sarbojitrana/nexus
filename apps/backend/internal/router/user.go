package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sarbojitrana/nexus/internal/handler"
	"github.com/sarbojitrana/nexus/internal/middleware"
)

// registerUserRoutes registers profile CRUD, search, and posts-by-user
// endpoints. Static segments (e.g. /users/me) take priority over the
// /users/:id param route in Echo's router regardless of registration order.
func registerUserRoutes(g *echo.Group, h *handler.Handlers, mw *middleware.Middlewares) {
	g.GET("/users", h.User.GetUsers, mw.Auth.OptionalAuth)
	g.GET("/users/:id", h.User.GetUserByID, mw.Auth.OptionalAuth)
	g.GET("/users/:id/posts", h.User.GetPostsByUserID, mw.Auth.OptionalAuth)

	g.GET("/users/me", h.User.GetMe, mw.Auth.RequireAuth)
	g.PATCH("/users/me", h.User.UpdateMe, mw.Auth.RequireAuth)
	g.DELETE("/users/me", h.User.DeleteMe, mw.Auth.RequireAuth)
}
