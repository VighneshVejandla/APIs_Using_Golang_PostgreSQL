package messaging

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

var rabbitmqConn *amqp.Connection

// ConnectRabbitMQ establishes a connection to RabbitMQ
func ConnectRabbitMQ() (*amqp.Connection, error) {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	log.Println("Successfully connected to RabbitMQ")
	return conn, nil
}
