package rabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rabbitURL = "amqp://dev:dev@localhost:5672/"
	exchange  = "events.fanout"
	exchType  = "fanout"
)

func (r *rb) Publish(queName string, eventBody []byte) error {
	conn, ch, err := r.connect(rabbitURL)
	defer func() {
		conn.Close()
		ch.Close()
	}()

	err = ch.ExchangeDeclare(
		exchange,
		exchType,
		true,  // durable
		false, // auto-delete
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		exchange,
		"", // routing key ignored for fanout
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
