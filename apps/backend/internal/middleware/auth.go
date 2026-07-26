package middleware

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/labstack/echo/v4"
	"github.com/sarbojitrana/nexus/internal/errs"
	"github.com/sarbojitrana/nexus/internal/server"
)

type AuthMiddleware struct {
	server *server.Server
}

func NewAuthMiddleware(s *server.Server) *AuthMiddleware {
	return &AuthMiddleware{
		server: s,
	}
}

// writeUnauthorized writes the standard 401 JSON error body. Shared between
// RequireAuth (missing/invalid session) and OptionalAuth (invalid session --
// a route using OptionalAuth still rejects garbage credentials, it just
// doesn't require credentials to be present at all).
func (auth *AuthMiddleware) writeUnauthorized(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	w.Header().Set("Content-Type", "application/json") // sets Content-Type header
	w.WriteHeader(http.StatusUnauthorized)             // sets status code

	response := map[string]string{ // response body
		"code":     "UNAUTHORIZED",
		"message":  "Unauthorized",
		"override": "false",
		"status":   "401",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil { // encode to the response format
		auth.server.Logger.Error().Err(err).Str("function", "RequireAuth").Dur("duration", time.Since(start)).Msg("failed to write JSON response")
	} else {
		auth.server.Logger.Error().Str("function", "RequireAuth").Dur("duration", time.Since(start)).Msg("could not get session claims from context")
	}
}

func (auth *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.WrapMiddleware( // echo middleware expects echo.HandlerFunc but clerk provides http handler so we need to wrap it
		clerkhttp.WithHeaderAuthorization( // does all the header authorization
			clerkhttp.AuthorizationFailureHandler(http.HandlerFunc(auth.writeUnauthorized)),
		))(func(c echo.Context) error {
		start := time.Now()

		claims, ok := clerk.SessionClaimsFromContext(c.Request().Context()) // take out claims from the jwt token from the request

		if !ok {
			auth.server.Logger.Error().
				Str("function", "RequireAuth").
				Str("request_id", GetRequestID(c)).
				Dur("duration", time.Since(start)).
				Msg("could not get session claims from context")
			return errs.NewUnauthorizedError("Unauthorized", false)
		}

		c.Set("user_id", claims.Subject)
		c.Set("user_role", claims.ActiveOrganizationRole)
		c.Set("permissions", claims.Claims.ActiveOrganizationPermissions)

		auth.server.Logger.Info().
			Str("function", "RequireAuth").
			Str("user_id", claims.Subject).
			Str("request_id", GetRequestID(c)).
			Dur("duration", time.Since(start)).
			Msg("user authenticated successfully")

		return next(c)

	})

}

// OptionalAuth extracts session claims into context when a valid Authorization
// header is present, but never requires one -- routes using this middleware
// serve anonymous requests normally. An invalid/expired token is still
// rejected with 401, since that's a genuine error rather than "no session".
func (auth *AuthMiddleware) OptionalAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.WrapMiddleware(
		clerkhttp.WithHeaderAuthorization(
			clerkhttp.AuthorizationFailureHandler(http.HandlerFunc(auth.writeUnauthorized)),
		))(func(c echo.Context) error {
		if claims, ok := clerk.SessionClaimsFromContext(c.Request().Context()); ok {
			c.Set("user_id", claims.Subject)
			c.Set("user_role", claims.ActiveOrganizationRole)
			c.Set("permissions", claims.Claims.ActiveOrganizationPermissions)
		}

		return next(c)
	})
}
