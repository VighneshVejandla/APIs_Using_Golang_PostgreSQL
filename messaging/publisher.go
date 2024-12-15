package messaging

import (
	"log"

	"github.com/streadway/amqp"
)

var rabbitMQChannel *amqp.Channel

// Initialize the RabbitMQ publisher
func InitializePublisher(amqpURL string) error {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	rabbitMQChannel = ch
	return nil
}

// Publish a message to a specific RabbitMQ queue
func PublishMessage(queueName, message string) error {
	if rabbitMQChannel == nil {
		return amqp.ErrClosed
	}

	err := rabbitMQChannel.Publish(
		"",       // exchange
		queueName, // routing key (queue name)
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Printf("Message published to queue %s: %s", queueName, message)
	return nil
}
