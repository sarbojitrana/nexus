package router

import (
	"github.com/labstack/echo/v4"
	"github.com/sarbojitrana/nexus/internal/handler"
)

// registerWebhookRoutes registers inbound webhook endpoints. These are
// authenticated by provider signature (not a Clerk session), so they are
// deliberately kept outside any RequireAuth-guarded route group.
func registerWebhookRoutes(r *echo.Echo, h *handler.Handlers) {
	r.POST("/webhooks/clerk", h.Webhook.HandleClerkWebhook)
}
