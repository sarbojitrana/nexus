package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sarbojitrana/nexus/internal/handler"
	"github.com/sarbojitrana/nexus/internal/middleware"
)

// registerFollowRoutes registers user-follow routes. Community-follow routes
// (also served by FollowHandler) are registered in registerCommunityRoutes,
// since they hang off the /communities/:id path.
func registerFollowRoutes(g *echo.Group, h *handler.Handlers, mw *middleware.Middlewares) {
	g.POST("/users/:id/follow", h.Follow.FollowUser, mw.Auth.RequireAuth)
	g.DELETE("/users/:id/follow", h.Follow.UnFollowUser, mw.Auth.RequireAuth)
	g.GET("/users/:id/follow", h.Follow.IsFollowingUser, mw.Auth.RequireAuth)
	g.GET("/users/:id/followers", h.Follow.GetFollowers, mw.Auth.OptionalAuth)
	g.GET("/users/:id/following", h.Follow.GetFollowing, mw.Auth.OptionalAuth)
}
