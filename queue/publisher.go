package queue

import (
	"context"
	"encoding/json"
	"time"

	"gin-project/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

// EmailJobType distinguishes which email a worker should send.
type EmailJobType string

const (
	WelcomeEmail       EmailJobType = "welcome"
	OTPEmail           EmailJobType = "otp"
	ResetPasswordEmail EmailJobType = "reset_password"
)

// EmailJob is the JSON message body published to the email queue.
type EmailJob struct {
	Type          EmailJobType `json:"type"`
	Email         string       `json:"email"`
	Name          string       `json:"name"`
	OTP           string       `json:"otp,omitempty"`
	ResetLink     string       `json:"resetLink,omitempty"`
	ExpiryMinutes int          `json:"expiryMinutes,omitempty"`
}

// publish marshals body to JSON and sends it to the named queue using the
// default exchange (routing key == queue name).
func publish(queueName string, body any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return config.RabbitChannel.PublishWithContext(
		ctx,
		"",        // exchange — "" is the default direct exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // persist the message to disk
			Body:         data,
		},
	)
}

// PublishWelcomeEmail enqueues a welcome-email job for a newly registered user.
func PublishWelcomeEmail(email, name string) error {
	return publish(EmailQueue, EmailJob{
		Type:  WelcomeEmail,
		Email: email,
		Name:  name,
	})
}

// PublishOTPEmail enqueues an OTP-email job.
func PublishOTPEmail(email, name, otp string, expiryMinutes int) error {
	return publish(EmailQueue, EmailJob{
		Type:          OTPEmail,
		Email:         email,
		Name:          name,
		OTP:           otp,
		ExpiryMinutes: expiryMinutes,
	})
}

// PublishResetPasswordEmail enqueues a password-reset email job.
func PublishResetPasswordEmail(email, name, resetLink string, expiryMinutes int) error {
	return publish(EmailQueue, EmailJob{
		Type:          ResetPasswordEmail,
		Email:         email,
		Name:          name,
		ResetLink:     resetLink,
		ExpiryMinutes: expiryMinutes,
	})
}
