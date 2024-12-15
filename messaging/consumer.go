package messaging

import (
	"log"

	"github.com/streadway/amqp"
)

func StartConsumer(amqpURL, queueName string) error {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received message: %s", msg.Body)
		}
	}()

	return nil
}
