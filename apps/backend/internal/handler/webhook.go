package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	svix "github.com/svix/svix-webhooks/go"

	"github.com/sarbojitrana/nexus/internal/errs"
	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/server"
	"github.com/sarbojitrana/nexus/internal/service"
)

const (
	clerkEventUserCreated    = "user.created"
	clerkEventUserUpdated    = "user.updated"
	clerkEventUserDeleted    = "user.deleted"
	clerkEventSessionCreated = "session.created"
)

type clerkWebhookEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ClerkWebhookHandler struct {
	Handler
	webhook  *svix.Webhook
	services *service.Services
}

func NewClerkWebhookHandler(s *server.Server, services *service.Services) *ClerkWebhookHandler {
	wh, err := svix.NewWebhook(s.Config.Auth.WebhookSecret)
	if err != nil {
		s.Logger.Fatal().Err(err).Msg("invalid clerk webhook secret")
	}

	return &ClerkWebhookHandler{
		Handler:  NewHandler(s),
		webhook:  wh,
		services: services,
	}
}

func (h *ClerkWebhookHandler) HandleClerkWebhook(c echo.Context) error {
	logger := middleware.GetLogger(c)

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		logger.Error().Err(err).Msg("failed to read clerk webhook body")
		return errs.NewBadRequestError("could not read request body", false, nil, nil, nil)
	}

	if err := h.webhook.Verify(body, c.Request().Header); err != nil {
		logger.Warn().Err(err).Msg("clerk webhook signature verification failed")
		return errs.NewBadRequestError("invalid webhook signature", false, nil, nil, nil)
	}

	var event clerkWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		logger.Error().Err(err).Msg("failed to parse clerk webhook payload")
		return errs.NewBadRequestError("malformed webhook payload", false, nil, nil, nil)
	}

	if txn := newrelic.FromContext(c.Request().Context()); txn != nil {
		txn.AddAttribute("webhook.event_type", event.Type)
	}

	ctx := c.Request().Context()

	switch event.Type {
	case clerkEventUserCreated:
		var cu clerk.User
		if err := json.Unmarshal(event.Data, &cu); err != nil {
			logger.Error().Err(err).Msg("failed to parse clerk user.created payload")
			return errs.NewBadRequestError("malformed webhook payload", false, nil, nil, nil)
		}

		if _, err := h.services.User.CreateFromClerk(ctx, &cu); err != nil {
			logger.Error().Err(err).Str("clerk_user_id", cu.ID).Msg("failed to sync created user")
			return err
		}

		logger.Info().Str("clerk_user_id", cu.ID).Msg("synced user.created from clerk")

	case clerkEventUserUpdated:
		var cu clerk.User
		if err := json.Unmarshal(event.Data, &cu); err != nil {
			logger.Error().Err(err).Msg("failed to parse clerk user.updated payload")
			return errs.NewBadRequestError("malformed webhook payload", false, nil, nil, nil)
		}

		if _, err := h.services.User.UpdateFromClerk(ctx, &cu); err != nil {
			logger.Error().Err(err).Str("clerk_user_id", cu.ID).Msg("failed to sync updated user")
			return err
		}

		logger.Info().Str("clerk_user_id", cu.ID).Msg("synced user.updated from clerk")

	case clerkEventUserDeleted:
		var payload struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(event.Data, &payload); err != nil {
			logger.Error().Err(err).Msg("failed to parse clerk user.deleted payload")
			return errs.NewBadRequestError("malformed webhook payload", false, nil, nil, nil)
		}

		if err := h.services.User.DeleteFromClerk(ctx, payload.ID); err != nil {
			logger.Error().Err(err).Str("clerk_user_id", payload.ID).Msg("failed to sync deleted user")
			return err
		}

		logger.Info().Str("clerk_user_id", payload.ID).Msg("synced user.deleted from clerk")

	case clerkEventSessionCreated:
		var sess clerk.Session
		if err := json.Unmarshal(event.Data, &sess); err != nil {
			logger.Error().Err(err).Msg("failed to parse clerk session.created payload")
			return errs.NewBadRequestError("malformed webhook payload", false, nil, nil, nil)
		}

		h.services.User.NotifySignIn(ctx, sess.UserID)

	default:
		logger.Debug().Str("event_type", event.Type).Msg("ignoring unhandled clerk webhook event")
	}

	return c.NoContent(http.StatusOK)
}
