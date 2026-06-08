package workers

import (
	"encoding/json"
	"fmt"
	"log"

	"gin-project/config"
	"gin-project/helper"
	"gin-project/queue"
)

// StartEmailWorker begins consuming the email queue and processes jobs in a
// background goroutine. Each message is acknowledged only after it is handled
// successfully, so an unexpected crash won't lose in-flight emails.
func StartEmailWorker() error {
	messages, err := config.RabbitChannel.Consume(
		queue.EmailQueue, // queue
		"",               // consumer tag (auto-generated)
		false,            // auto-ack — false means we ack manually
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			var job queue.EmailJob

			if err := json.Unmarshal(msg.Body, &job); err != nil {
				log.Printf("[email-worker] skipping malformed job: %v", err)
				_ = msg.Nack(false, false) // don't requeue garbage
				continue
			}

			fmt.Printf("Received email job: %+v\n", job)

			if err := handleEmailJob(job); err != nil {
				log.Printf("[email-worker] failed to send %s email to %s: %v", job.Type, job.Email, err)
				// Don't requeue: a failing email would loop forever and hammer
				// SMTP. Add a dead-letter queue later for retries.
				_ = msg.Nack(false, false)
				continue
			}

			log.Printf("[email-worker] sent %s email to %s", job.Type, job.Email)
			_ = msg.Ack(false)
		}
	}()

	log.Println("Email worker started")
	return nil
}

// handleEmailJob routes a job to the right email helper based on its type.
func handleEmailJob(job queue.EmailJob) error {
	switch job.Type {
	case queue.WelcomeEmail:
		return helper.SendWelcomeEmail(job.Email, job.Name)
	case queue.OTPEmail:
		return helper.SendOTPEmail(job.Email, job.Name, job.OTP, job.ExpiryMinutes)
	default:
		log.Printf("[email-worker] unknown job type %q, ignoring", job.Type)
		return nil
	}
}
