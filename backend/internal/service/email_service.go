package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/nats-io/nats.go"

	"github.com/bifshteksex/hertzboard/internal/config"
)

// EmailService handles email sending
type EmailService struct {
	cfg  *config.EmailConfig
	nats *nats.Conn
}

// EmailMessage represents an email message
type EmailMessage struct {
	To      string                 `json:"to"`
	Subject string                 `json:"subject"`
	Type    string                 `json:"type"`
	Data    map[string]interface{} `json:"data"`
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.EmailConfig, nc *nats.Conn) *EmailService {
	return &EmailService{
		cfg:  cfg,
		nats: nc,
	}
}

// PublishEmail publishes an email message to NATS queue
func (s *EmailService) PublishEmail(msg *EmailMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal email message: %w", err)
	}

	if err := s.nats.Publish("emails", data); err != nil {
		return fmt.Errorf("failed to publish email: %w", err)
	}

	return nil
}

// SendWelcomeEmail sends a welcome email
func (s *EmailService) SendWelcomeEmail(to, name string) error {
	return s.PublishEmail(&EmailMessage{
		To:      to,
		Subject: "Welcome to HertzBoard!",
		Type:    "welcome",
		Data: map[string]interface{}{
			"name": name,
		},
	})
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(to, name, token, resetURL string) error {
	return s.PublishEmail(&EmailMessage{
		To:      to,
		Subject: "Reset your password",
		Type:    "password_reset",
		Data: map[string]interface{}{
			"name":      name,
			"token":     token,
			"reset_url": resetURL,
		},
	})
}

// SendEmailVerification sends an email verification
func (s *EmailService) SendEmailVerification(to, name, token, verifyURL string) error {
	return s.PublishEmail(&EmailMessage{
		To:      to,
		Subject: "Verify your email",
		Type:    "email_verification",
		Data: map[string]interface{}{
			"name":       name,
			"token":      token,
			"verify_url": verifyURL,
		},
	})
}

// SendWorkspaceInvite sends a workspace invitation email
func (s *EmailService) SendWorkspaceInvite(to, workspaceName, inviterName, inviteURL string) error {
	return s.PublishEmail(&EmailMessage{
		To:      to,
		Subject: fmt.Sprintf("You've been invited to %s", workspaceName),
		Type:    "workspace_invite",
		Data: map[string]interface{}{
			"workspace_name": workspaceName,
			"inviter_name":   inviterName,
			"invite_url":     inviteURL,
		},
	})
}

// EmailWorker processes email messages from NATS queue
type EmailWorker struct {
	cfg  *config.EmailConfig
	nats *nats.Conn
	sub  *nats.Subscription
}

// NewEmailWorker creates a new email worker
func NewEmailWorker(cfg *config.EmailConfig, nc *nats.Conn) (*EmailWorker, error) {
	worker := &EmailWorker{
		cfg:  cfg,
		nats: nc,
	}

	// Subscribe to email queue
	sub, err := nc.QueueSubscribe("emails", "email-workers", worker.handleMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to email queue: %w", err)
	}

	worker.sub = sub
	return worker, nil
}

// Close closes the email worker subscription
func (w *EmailWorker) Close() error {
	if w.sub != nil {
		return w.sub.Unsubscribe()
	}
	return nil
}

// handleMessage processes an email message
func (w *EmailWorker) handleMessage(msg *nats.Msg) {
	var emailMsg EmailMessage
	if err := json.Unmarshal(msg.Data, &emailMsg); err != nil {
		fmt.Printf("Failed to unmarshal email message: %v\n", err)
		return
	}

	if err := w.sendEmail(&emailMsg); err != nil {
		fmt.Printf("Failed to send email to %s: %v\n", emailMsg.To, err)
		// TODO: Implement retry logic with exponential backoff
		return
	}

	fmt.Printf("Email sent successfully to %s\n", emailMsg.To)
}

// sendEmail sends an actual email via SMTP
func (w *EmailWorker) sendEmail(msg *EmailMessage) error {
	// Generate email body from template
	body, err := w.renderTemplate(msg.Type, msg.Data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Prepare email
	from := w.cfg.From
	to := msg.To
	subject := msg.Subject

	message := fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body

	// Send via SMTP
	addr := fmt.Sprintf("%s:%d", w.cfg.SMTPHost, w.cfg.SMTPPort)

	// For development (MailHog), we don't need authentication
	if w.cfg.SMTPUser == "" && w.cfg.SMTPPassword == "" {
		// Connect without auth
		c, err := smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer c.Close()

		if err := c.Mail(from); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		if err := c.Rcpt(to); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}

		wc, err := c.Data()
		if err != nil {
			return fmt.Errorf("failed to create data writer: %w", err)
		}
		defer wc.Close()

		if _, err := wc.Write([]byte(message)); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

		return nil
	}

	// For production with authentication
	auth := smtp.PlainAuth("", w.cfg.SMTPUser, w.cfg.SMTPPassword, w.cfg.SMTPHost)
	err = smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// renderTemplate renders an email template
func (w *EmailWorker) renderTemplate(templateType string, data map[string]interface{}) (string, error) {
	templates := map[string]string{
		"welcome": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body>
    <h1>Welcome to HertzBoard, {{.name}}!</h1>
    <p>We're excited to have you on board.</p>
    <p>Get started by creating your first workspace and start collaborating!</p>
</body>
</html>
`,
		"password_reset": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body>
    <h1>Reset your password</h1>
    <p>Hello {{.name}},</p>
    <p>You requested to reset your password. Click the link below to continue:</p>
    <p><a href="{{.reset_url}}?token={{.token}}">Reset Password</a></p>
    <p>This link will expire in 1 hour.</p>
    <p>If you didn't request this, you can safely ignore this email.</p>
</body>
</html>
`,
		"email_verification": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body>
    <h1>Verify your email</h1>
    <p>Hello {{.name}},</p>
    <p>Please verify your email address by clicking the link below:</p>
    <p><a href="{{.verify_url}}?token={{.token}}">Verify Email</a></p>
</body>
</html>
`,
		"workspace_invite": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body>
    <h1>You've been invited to {{.workspace_name}}</h1>
    <p>{{.inviter_name}} has invited you to collaborate on {{.workspace_name}}.</p>
    <p><a href="{{.invite_url}}">Accept Invitation</a></p>
</body>
</html>
`,
	}

	tmplStr, exists := templates[templateType]
	if !exists {
		return "", fmt.Errorf("template not found: %s", templateType)
	}

	tmpl, err := template.New(templateType).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
