package email

import (
	"context"
	"fmt"
	"log"
	"muchway/event_service/domain"
	"net/smtp"
	"time"
)

// EmailService defines the interface for sending emails
type EmailService interface {
	SendNewEventNotification(ctx context.Context, users []string, event *domain.Event) error
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

// SendNewEventNotification sends a notification email about a new event to all users
func (s *emailService) SendNewEventNotification(ctx context.Context, users []string, event *domain.Event) error {
	if len(users) == 0 {
		return fmt.Errorf("no users to send notification to")
	}

	// SMTP server configuration
	smtpHost := s.config.SMTPHost
	smtpPort := s.config.SMTPPort
	smtpAddr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// Sender credentials
	from := s.config.SenderEmail
	password := s.config.SenderPasswd

	// Format the start time
	formattedTime := event.StartTime.Format("Monday, January 2, 2006 at 3:04 PM")

	// Email content
	subject := fmt.Sprintf("New Event: %s - MuchWayBet", event.Name)
	body := fmt.Sprintf(`
	<html>
	<body>
		<h2>New Event: %s</h2>
		<p>We're excited to announce a new event on MuchWayBet!</p>
		<p><strong>Event Details:</strong></p>
		<ul>
			<li><strong>Name:</strong> %s</li>
			<li><strong>Start Time:</strong> %s</li>
			<li><strong>Status:</strong> %s</li>
		</ul>
		<p>Don't miss out on this opportunity to place your bets!</p>
		<p>Visit MuchWayBet now to participate.</p>
		<p>Best regards,<br>The MuchWayBet Team</p>
	</body>
	</html>
	`, event.Name, event.Name, formattedTime, event.Status)

	// Compose message
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := fmt.Sprintf("Subject: %s\n%s\n%s", subject, mime, body)

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send email to all users
	// To avoid exposing all email addresses, we'll send individual emails
	for _, to := range users {
		err := smtp.SendMail(smtpAddr, auth, from, []string{to}, []byte(message))
		if err != nil {
			log.Printf("Failed to send email to %s: %v", to, err)
			// Continue sending to other users even if one fails
		} else {
			log.Printf("Event notification email sent to: %s", to)
		}

		// Add a small delay to avoid overwhelming the SMTP server
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}
