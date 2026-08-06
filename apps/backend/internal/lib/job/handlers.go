package job

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/sarbojitrana/nexus/internal/config"
	"github.com/sarbojitrana/nexus/internal/lib/email"
)

var emailClient *email.Client

func (j *JobService) InitHandlers(config *config.Config, logger *zerolog.Logger) {
	emailClient = email.NewClient(config, logger)
}

func (j *JobService) handleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var p WeclomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmmarshal welcome email payload: %w", err)
	}

	err := emailClient.SendWelcomeEmail(
		p.To,
		p.FirstName,
	)

	if err != nil {
		j.logger.Error().
			Str("type", "welcome").
			Str("to", p.To).
			Err(err).
			Msg("Failed to send welcome email")
		return err
	}

	j.logger.Info().
		Str("type", "welcome").
		Str("to", p.To).
		Msg("Successfully sent welcome email")

	return nil
}

func (j *JobService) handleSignInEmailTask(ctx context.Context, t *asynq.Task) error {
	var p SignInEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal signin email payload: %w", err)
	}

	err := emailClient.SendSignInEmail(p.To, p.FirstName, p.SignInTime)
	if err != nil {
		j.logger.Error().
			Str("type", "signin").
			Str("to", p.To).
			Err(err).
			Msg("Failed to send sign-in email")
		return err
	}

	j.logger.Info().
		Str("type", "signin").
		Str("to", p.To).
		Msg("Successfully sent sign-in email")

	return nil
}

func (j *JobService) handlePasswordChangedEmailTask(ctx context.Context, t *asynq.Task) error {
	var p PasswordChangedEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal password changed email payload: %w", err)
	}

	err := emailClient.SendPasswordChangedEmail(p.To, p.FirstName, p.ChangedAt)
	if err != nil {
		j.logger.Error().
			Str("type", "password_changed").
			Str("to", p.To).
			Err(err).
			Msg("Failed to send password changed email")
		return err
	}

	j.logger.Info().
		Str("type", "password_changed").
		Str("to", p.To).
		Msg("Successfully sent password changed email")

	return nil
}
