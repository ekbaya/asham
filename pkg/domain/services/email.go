package services

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
	"time"
)

// EmailConfig holds the configuration for the email service
type EmailConfig struct {
	Host              string
	Port              int
	Username          string
	Password          string
	From              string
	EmailTemplatePath string
}

// EmailService handles sending emails
type EmailService struct {
	config *EmailConfig
}

// NewEmailService creates a new email service
func NewEmailService(config *EmailConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

// SendWelcomeEmail sends a styled HTML welcome email to a new user with their password
func (s *EmailService) SendWelcomeEmail(toEmail, name, password string) error {
	subject := "Welcome to Our Service"

	// Read the email template from file
	tmpl, err := os.ReadFile(s.config.EmailTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read email template: %w", err)
	}

	// Parse the template
	t, err := template.New("welcomeEmail").Parse(string(tmpl))
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Prepare template data including band name "ASHAM" and dynamic year
	data := struct {
		Name     string
		Password string
		BandName string
		Year     int
	}{
		Name:     name,
		Password: password,
		BandName: "ASHAM",
		Year:     time.Now().Year(),
	}

	// Execute template with data
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	return s.sendEmail(toEmail, subject, body.String())
}

// sendEmail handles the actual email sending
func (s *EmailService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	from := s.config.From
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", from, to, subject, body)

	// Send the email
	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
	if err != nil {
		return errors.New("failed to send email: " + err.Error())
	}

	return nil
}
