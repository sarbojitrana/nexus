package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sarbojitrana/nexus/internal/handler"
	"github.com/sarbojitrana/nexus/internal/middleware"
)

func registerCommunityRoutes(g *echo.Group, h *handler.Handlers, mw *middleware.Middlewares) {
	g.POST("/communities", h.Community.CreateCommunity, mw.Auth.RequireAuth)
	g.GET("/communities", h.Community.GetCommunities)
	g.GET("/communities/:idOrSlug", h.Community.GetCommunityByIDOrSlug, mw.Auth.OptionalAuth)

	g.PATCH("/communities/:id", h.Community.UpdateSettings, mw.Auth.RequireAuth)
	g.DELETE("/communities/:id", h.Community.DeleteCommunity, mw.Auth.RequireAuth)

	g.POST("/communities/:id/join", h.Community.JoinCommunity, mw.Auth.RequireAuth)
	g.POST("/communities/:id/leave", h.Community.LeaveCommunity, mw.Auth.RequireAuth)

	g.POST("/communities/:id/follow", h.Follow.FollowCommunity, mw.Auth.RequireAuth)
	g.DELETE("/communities/:id/follow", h.Follow.UnFollowCommunity, mw.Auth.RequireAuth)
	g.GET("/communities/:id/follow", h.Follow.IsFollowingCommunity, mw.Auth.RequireAuth)

	g.GET("/communities/:id/members", h.Community.GetMembers, mw.Auth.OptionalAuth)
	g.PATCH("/communities/:id/members/:userId/role", h.Community.ChangeMemberRole, mw.Auth.RequireAuth)

	g.GET("/communities/:id/posts/:postId", h.Community.GetCommunityPostByID, mw.Auth.OptionalAuth)
	g.DELETE("/communities/:id/posts/:postId", h.Community.DeleteCommunityPost, mw.Auth.RequireAuth)

	g.POST("/communities/:id/bans", h.Community.BanMember, mw.Auth.RequireAuth)

	g.POST("/communities/:id/reports", h.Community.ReportPost, mw.Auth.RequireAuth)
	g.GET("/communities/:id/reports", h.Community.GetReports, mw.Auth.RequireAuth)
	g.GET("/communities/:id/reports/:reportId", h.Community.GetReportByID, mw.Auth.RequireAuth)
	g.PATCH("/communities/:id/reports/:reportId", h.Community.ResolveReport, mw.Auth.RequireAuth)
}
