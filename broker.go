package main

import (
	"context"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
    Host string
    Port string
    User string
    Pass string
}

func GetChannel() *amqp.Channel {
    config := ParseConfig()
    url := fmt.Sprintf("amqp://%s:%s@%s:%s", config.User, config.Pass, config.Host, config.Port)
    conn, err := amqp.Dial(url)
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    _, err = ch.QueueDeclare(
        "events", // name
        false,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    failOnError(err, "Failed to declare a queue")

    return ch
}

func PublishEventMessage(eventId string) {
    ch := GetChannel()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    body := fmt.Sprintf(`{"event": "%s"}`, eventId)
    err := ch.PublishWithContext(ctx,
    "",     // exchange
    "events", // routing key Should be same as queue name
    false,  // mandatory
    false,  // immediate
    amqp.Publishing {
        ContentType: "application/json",
        Body:        []byte(body),
    })
    failOnError(err, "Failed to publish a message")
}

func ParseConfig() *RabbitMQConfig {
    return &RabbitMQConfig{
        Host: os.Getenv("RABBITMQ_HOST"),
        Port: os.Getenv("RABBITMQ_PORT"),
        User: os.Getenv("RABBITMQ_USER"),
        Pass: os.Getenv("RABBITMQ_PASS"),
    }
}
