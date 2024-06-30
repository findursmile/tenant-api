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

func GetChannel() (*amqp.Channel, error) {
    config := ParseConfig()
    url := fmt.Sprintf("amqp://%s:%s@%s:%s", config.User, config.Pass, config.Host, config.Port)

    conn, err := amqp.Dial(url)

    if err != nil {
        failOnError(err, "Failed to open a channel")
        return nil, err
    }

    defer conn.Close()

    ch, err := conn.Channel()

    if err != nil {
        failOnError(err, "Failed to open a channel")
        return nil, err
    }

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

    return ch, nil
}

func PublishEventMessage(eventId string) {
    ch, err := GetChannel()

    if err != nil {
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    body := fmt.Sprintf(`{"event": "%s"}`, eventId)
    err = ch.PublishWithContext(ctx,
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
