package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sarbojitrana/nexus/internal/handler"
	"github.com/sarbojitrana/nexus/internal/middleware"
)

func registerPostRoutes(g *echo.Group, h *handler.Handlers, mw *middleware.Middlewares) {
	g.POST("/posts", h.Post.CreatePost, mw.Auth.RequireAuth)
	g.PATCH("/posts/:id", h.Post.UpdatePost, mw.Auth.RequireAuth)
	g.DELETE("/posts/:id", h.Post.DeletePost, mw.Auth.RequireAuth)
	g.GET("/posts/:id", h.Post.GetPostByID, mw.Auth.OptionalAuth)
	g.GET("/posts/:id/comments", h.Post.GetComments, mw.Auth.OptionalAuth)
	g.GET("/posts/:id/replies", h.Post.GetReplies, mw.Auth.OptionalAuth)
	g.POST("/posts/:id/react", h.Post.React, mw.Auth.RequireAuth)

	g.GET("/feed", h.Post.GetFeed, mw.Auth.OptionalAuth)
}
