package email

import (
	"fmt"
	"net/smtp"
)

// EmailService defines the interface for sending emails
type EmailService interface {
	SendRegistrationConfirmation(to, username string) error
}

// Config holds email configuration
type Config struct {
	SMTPHost     string
	SMTPPort     string
	SenderEmail  string
	SenderPasswd string
}

// emailService implements EmailService
type emailService struct {
	config Config
}

// NewEmailService creates a new email service
func NewEmailService(config Config) EmailService {
	return &emailService{
		config: config,
	}
}

// SendRegistrationConfirmation sends a confirmation email to a newly registered user
func (s *emailService) SendRegistrationConfirmation(to, username string) error {
	// SMTP server configuration
	smtpHost := s.config.SMTPHost
	smtpPort := s.config.SMTPPort
	smtpAddr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// Sender credentials
	from := s.config.SenderEmail
	password := s.config.SenderPasswd

	// Email content
	subject := "Welcome to MuchWayBet - Registration Successful!"
	body := fmt.Sprintf(`
	<html>
	<body>
		<h2>Welcome to MuchWayBet, %s!</h2>
		<p>Your registration was successful. Thank you for joining our platform!</p>
		<p>You can now log in and start using all the features of MuchWayBet.</p>
		<p>If you have any questions or need assistance, please don't hesitate to contact our support team.</p>
		<p>Best regards,<br>The MuchWayBet Team</p>
	</body>
	</html>
	`, username)

	// Compose message
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := fmt.Sprintf("Subject: %s\n%s\n%s", subject, mime, body)

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send email
	err := smtp.SendMail(smtpAddr, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
