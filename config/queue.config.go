package config

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

func InitRabbitMQ() {
	var err error

	RabbitConn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("Failed to connect RabbitMQ: %v", err)
	}

	RabbitChannel, err = RabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}

	log.Println("RabbitMQ Connected")
}

func CloseRabbitMQ() {
	if RabbitChannel != nil {
		RabbitChannel.Close()
	}

	if RabbitConn != nil {
		RabbitConn.Close()
	}
}
