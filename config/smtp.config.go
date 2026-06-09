package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

// SendMail sends an HTML email using the configured transport.
//
// Transport is chosen by MAIL_TRANSPORT:
//   - "brevo" (or "http"): send via the Brevo HTTP API over port 443. Use this
//     in production on hosts that block outbound SMTP ports (e.g. Render blocks
//     25/465/587).
//   - anything else / unset: send via SMTP (good for local development).
func SendMail(to, subject, body string) error {
	switch os.Getenv("MAIL_TRANSPORT") {
	case "brevo", "http":
		return sendViaBrevo(to, subject, body)
	default:
		return sendViaSMTP(to, subject, body)
	}
}

// mailFrom returns the sender address (MAIL_FROM, falling back to SMTP_USER).
func mailFrom() string {
	if from := os.Getenv("MAIL_FROM"); from != "" {
		return from
	}
	return os.Getenv("SMTP_USER")
}

// mailFromName returns the sender display name.
func mailFromName() string {
	if name := os.Getenv("MAIL_FROM_NAME"); name != "" {
		return name
	}
	if name := os.Getenv("APP_NAME"); name != "" {
		return name
	}
	return "MyAyur"
}

func sendViaSMTP(to, subject, body string) error {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return fmt.Errorf("invalid SMTP_PORT: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", mailFrom())
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)

	return d.DialAndSend(m)
}

// brevoEmailRequest is the JSON body for POST https://api.brevo.com/v3/smtp/email
type brevoEmailRequest struct {
	Sender      brevoContact   `json:"sender"`
	To          []brevoContact `json:"to"`
	Subject     string         `json:"subject"`
	HTMLContent string         `json:"htmlContent"`
}

type brevoContact struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email"`
}

func sendViaBrevo(to, subject, body string) error {
	apiKey := os.Getenv("BREVO_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("BREVO_API_KEY is not set")
	}

	from := mailFrom()
	if from == "" {
		return fmt.Errorf("MAIL_FROM (or SMTP_USER) is not set")
	}

	payload := brevoEmailRequest{
		Sender:      brevoContact{Name: mailFromName(), Email: from},
		To:          []brevoContact{{Email: to}},
		Subject:     subject,
		HTMLContent: body,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.brevo.com/v3/smtp/email", bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("api-key", apiKey)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("brevo API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}
