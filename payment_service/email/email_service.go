package email

import (
	"fmt"
	"net/smtp"
)

// EmailService defines the interface for sending emails
type EmailService interface {
	SendPaymentConfirmation(to, userID string, amount float64, paymentType, status string) error
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

// SendPaymentConfirmation sends a confirmation email for a payment
func (s *emailService) SendPaymentConfirmation(to, userID string, amount float64, paymentType, status string) error {
	// SMTP server configuration
	smtpHost := s.config.SMTPHost
	smtpPort := s.config.SMTPPort
	smtpAddr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// Sender credentials
	from := s.config.SenderEmail
	password := s.config.SenderPasswd

	// Email content
	var subject, bodyTemplate string

	if paymentType == "deposit" {
		subject = "MuchWayBet - Deposit Confirmation"
		bodyTemplate = `
		<html>
		<body>
			<h2>Deposit Confirmation</h2>
			<p>Dear MuchWayBet User,</p>
			<p>Your deposit of <strong>%.2f</strong> has been processed successfully.</p>
			<p>Transaction Status: <strong>%s</strong></p>
			<p>Thank you for using MuchWayBet!</p>
			<p>Best regards,<br>The MuchWayBet Team</p>
		</body>
		</html>
		`
	} else if paymentType == "withdraw" {
		subject = "MuchWayBet - Withdrawal Confirmation"
		bodyTemplate = `
		<html>
		<body>
			<h2>Withdrawal Confirmation</h2>
			<p>Dear MuchWayBet User,</p>
			<p>Your withdrawal of <strong>%.2f</strong> has been processed.</p>
			<p>Transaction Status: <strong>%s</strong></p>
			<p>Thank you for using MuchWayBet!</p>
			<p>Best regards,<br>The MuchWayBet Team</p>
		</body>
		</html>
		`
	} else {
		subject = "MuchWayBet - Payment Notification"
		bodyTemplate = `
		<html>
		<body>
			<h2>Payment Notification</h2>
			<p>Dear MuchWayBet User,</p>
			<p>A payment of <strong>%.2f</strong> has been processed for your account.</p>
			<p>Payment Type: <strong>%s</strong></p>
			<p>Transaction Status: <strong>%s</strong></p>
			<p>Thank you for using MuchWayBet!</p>
			<p>Best regards,<br>The MuchWayBet Team</p>
		</body>
		</html>
		`
	}

	// Format the body with the payment details
	var body string
	if paymentType == "deposit" || paymentType == "withdraw" {
		body = fmt.Sprintf(bodyTemplate, amount, status)
	} else {
		body = fmt.Sprintf(bodyTemplate, amount, paymentType, status)
	}

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
