// messaging/initialize.go
package messaging

import (
    "github.com/streadway/amqp"
    "log"
)

// InitializeMessaging initializes the RabbitMQ connection
func InitializeMessaging(conn *amqp.Connection) {
    // Setup RabbitMQ channel, exchange, or other related setup
    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open RabbitMQ channel: %v", err)
    }
    defer ch.Close()

    log.Println("RabbitMQ messaging initialized")
}
