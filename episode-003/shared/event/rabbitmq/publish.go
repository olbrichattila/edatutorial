package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *rb) Publish(topic string, eventBody []byte) error {
	conn, ch, err := r.connect()
	if err != nil {
		return err
	}
	defer func() {
		if ch != nil {
			if closeErr := ch.Close(); closeErr != nil {
				fmt.Printf("Error closing channel: %v\n", closeErr)
			}
		}
		if conn != nil {
			if closeErr := conn.Close(); closeErr != nil {
				fmt.Printf("Error closing connection: %v\n", closeErr)
			}
		}
	}()

	err = r.declareOrCreateExchange(ch, topic)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		topic, // exchange name
		"",    // routing key ignored for fanout
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        eventBody,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
