package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/rs/zerolog"

	"github.com/sarbojitrana/nexus/internal/config"
)

type Client struct {
	addr     string
	from     string
	username string
	auth     smtp.Auth
	logger   *zerolog.Logger
}

func NewClient(cfg *config.Config, logger *zerolog.Logger) *Client {
	host := cfg.Integration.SMTPHost

	return &Client{
		addr:     host + ":" + cfg.Integration.SMTPPort,
		from:     cfg.Integration.SMTPFrom,
		username: cfg.Integration.SMTPUsername,
		auth:     smtp.PlainAuth("", cfg.Integration.SMTPUsername, cfg.Integration.SMTPPassword, host),
		logger:   logger,
	}
}

func (c *Client) SendEmail(to, subject string, templateName Template, data map[string]string) error {
	tmplPath := fmt.Sprintf("%s/%s.html", "templates/emails", templateName)

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse mail template %s: %w", templateName, err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute mail template %s: %w", templateName, err)
	}

	msg := mimeMessage(c.from, to, subject, body.String())

	if err := smtp.SendMail(c.addr, c.auth, c.username, []string{to}, msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func mimeMessage(from, to, subject, htmlBody string) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "From: %s\r\n", from)
	fmt.Fprintf(&buf, "To: %s\r\n", to)
	fmt.Fprintf(&buf, "Subject: %s\r\n", subject)
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(htmlBody)
	return buf.Bytes()
}
