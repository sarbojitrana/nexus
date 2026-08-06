package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sarbojitrana/nexus/internal/handler"
	"github.com/sarbojitrana/nexus/internal/middleware"
)

func registerSearchRoutes(g *echo.Group, h *handler.Handlers, mw *middleware.Middlewares) {
	g.GET("/search", h.Search.Search, mw.Auth.OptionalAuth)
}
