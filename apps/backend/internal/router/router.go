package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/sarbojitrana/nexus/internal/handler"
	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/server"
	"github.com/sarbojitrana/nexus/internal/service"
	"golang.org/x/time/rate"
)

func NewRouter(s *server.Server, h *handler.Handlers, services *service.Services) *echo.Echo {
	middlewares := middleware.NewMiddlewares(s)

	router := echo.New()

	router.HTTPErrorHandler = middlewares.Global.GlobalErrorHandler

	router.Use(
		echoMiddleware.RateLimiterWithConfig(echoMiddleware.RateLimiterConfig{
			Store: echoMiddleware.NewRateLimiterMemoryStore(rate.Limit(20)),
			DenyHandler: func(c echo.Context, identifier string, err error) error {
				if rateLimitMiddleware := middlewares.RateLimit; rateLimitMiddleware != nil {
					rateLimitMiddleware.RecordRateLimitHit(c.Path())
				}

				s.Logger.Warn().
					Str("request_id", middleware.GetRequestID(c)).
					Str("identifier", identifier).
					Str("path", c.Path()).
					Str("ip", c.RealIP()).
					Msg("rate limit exceeded")

				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			},
		}),

		middlewares.Global.CORS(),
		middlewares.Global.Secure(),
		middleware.RequestID(),
		middlewares.Tracing.NewRelicMiddleware(),
		middlewares.Tracing.EchanceTracing(),
		middlewares.ContextEnhancer.EnhanceContext(),
		middlewares.Global.RequestLogger(),
		middlewares.Global.Recover(),
	)

	registerSystemRoutes(router, h)

	registerWebhookRoutes(router, h)

	registerChatWebSocketRoute(router, h)

	v1 := router.Group("/api/v1")
	registerUserRoutes(v1, h, middlewares)
	registerPostRoutes(v1, h, middlewares)
	registerCommunityRoutes(v1, h, middlewares)
	registerFollowRoutes(v1, h, middlewares)
	registerChatRoutes(v1, h, middlewares)
	registerNotificationRoutes(v1, h, middlewares)
	registerSearchRoutes(v1, h, middlewares)
	registerStorageRoutes(v1, h, middlewares)

	return router

}
