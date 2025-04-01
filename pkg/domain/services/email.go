package services

import (
	"errors"
	"fmt"
	"net/smtp"
)

// EmailConfig holds the configuration for the email service
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// EmailService handles sending emails
type EmailService struct {
	config EmailConfig
}

// NewEmailService creates a new email service
func NewEmailService(config EmailConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

// SendWelcomeEmail sends a welcome email to a new user with their password
func (s *EmailService) SendWelcomeEmail(toEmail, name, password string) error {
	subject := "Welcome to Our Service"
	body := fmt.Sprintf(`
Hello %s,

Welcome to our service! Your account has been created successfully.

Your temporary password is: %s

Please login and change your password as soon as possible.

Best regards,
The Team
`, name, password)

	return s.sendEmail(toEmail, subject, body)
}

// sendEmail handles the actual email sending
func (s *EmailService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	from := s.config.From
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// Format the email
	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", from, to, subject, body)

	// Send the email
	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
	if err != nil {
		return errors.New("failed to send email: " + err.Error())
	}

	return nil
}
