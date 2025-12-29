package rabbitmq

import (
	"log"

	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *rb) Consume(topic, consumerName string, callback contracts.CallbackFn) error {
	conn, ch, err := r.connect(rabbitURL)
	defer func() {
		conn.Close()
		ch.Close()
	}()

	// Ensure exchange exists
	err = r.declareOrCreateExchange(ch, topic)
	if err != nil {
		return err
	}

	q, err := r.declareQueueIfNotDeclared(ch, topic, consumerName)
	if err != nil {
		return err
	}

	err = r.bindQueueIfNotDoneAlready(ch, topic, q.Name)
	if err != nil {
		return err
	}

	msgs, err := r.startConsume(ch, q.Name, consumerName)
	if err != nil {
		return err
	}

	log.Printf("consumer %s waiting for messages", consumerName)

	return r.consume(msgs, callback)
}

func (r *rb) declareQueueIfNotDeclared(ch *amqp.Channel, topic, consumerName string) (amqp.Queue, error) {
	queueName := topic + "." + consumerName
	return ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // Exclusive
		false, // NoWait
		nil,   // optional args
	)
}

func (r *rb) bindQueueIfNotDoneAlready(ch *amqp.Channel, topic, queueName string) error {
	return ch.QueueBind(
		queueName,
		"",    // Key ignored
		topic, // Exchange name
		false, // noWait
		nil,   // Args
	)
}

func (r *rb) startConsume(ch *amqp.Channel, queueName, consumerName string) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		queueName,
		consumerName,
		false, // manual ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // optional args
	)
}

func (r *rb) consume(msgs <-chan amqp.Delivery, callback contracts.CallbackFn) error {
	for msg := range msgs {
		err := callback(r, msg.Body)
		if err != nil {
			_ = msg.Nack(false, true) // Requeue
			continue
		}

		msg.Ack(false)
	}

	return nil
}
