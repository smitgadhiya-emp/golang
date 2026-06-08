package queue

import (
	"gin-project/config"
)

// Queue names used across the app.
const (
	EmailQueue = "email_queue"
)

// DeclareQueues ensures every queue exists before we publish to or consume from
// it. It is idempotent and should be called once at startup, after RabbitMQ is
// connected. (This is the Go equivalent of `assertQueue` in amqplib.)
func DeclareQueues() error {
	queues := []string{EmailQueue}

	for _, name := range queues {
		_, err := config.RabbitChannel.QueueDeclare(
			name,  // name
			true,  // durable — survives a broker restart
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return err
		}
	}

	return nil
}
