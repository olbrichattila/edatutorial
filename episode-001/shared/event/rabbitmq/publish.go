package rabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rabbitURL = "amqp://dev:dev@localhost:5672/"
	exchType  = "fanout"
)

func (r *rb) Publish(topic string, eventBody []byte) error {
	conn, ch, err := r.connect(rabbitURL)
	defer func() {
		conn.Close()
		ch.Close()
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
