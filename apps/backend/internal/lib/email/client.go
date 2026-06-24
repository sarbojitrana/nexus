package email

import (
	"github.com/rs/zerolog"

	"bytes"
	"fmt"
	"html/template"

	"github.com/pkg/errors"
	"github.com/resend/resend-go/v2"
	"github.com/sarbojitrana/nexus/internal/config"
)

type Client struct {
	client *resend.Client
	logger *zerolog.Logger
}

func NewClient(cfg *config.Config, logger *zerolog.Logger) *Client {
	return &Client{
		client: resend.NewClient(cfg.Integration.ResendAPIKey),
		logger: logger,
	}
}

func (c *Client) SendEmail(to, subject string, templateName Template, data map[string]string) error {
	tmplPath := fmt.Sprintf("%s/%s.html", "templates/emails", templateName)

	tmpl, err := template.ParseFiles(tmplPath)

	if err != nil {
		return errors.Wrapf(err, "failed to parse mail template %s", templateName)
	}

	var body bytes.Buffer

	if err := tmpl.Execute(&body, data); err != nil { // data is for placeholders in the mail
		return errors.Wrapf(err, "failed to execute mail template %s", templateName)
	}

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", "Nexus", "onboarding@resend.dev"),
		To:      []string{to},
		Subject: subject,
		Html:    body.String(),
	}

	_, err = c.client.Emails.Send(params)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
