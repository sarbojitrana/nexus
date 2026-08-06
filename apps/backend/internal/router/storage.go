package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sarbojitrana/nexus/internal/handler"
	"github.com/sarbojitrana/nexus/internal/middleware"
)

func registerStorageRoutes(g *echo.Group, h *handler.Handlers, mw *middleware.Middlewares) {
	g.POST("/uploads/presign", h.Storage.PresignUpload, mw.Auth.RequireAuth)
}
