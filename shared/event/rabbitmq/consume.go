package rabbitmq

import (
	"log"
)

func (r *rb) Consume(consumerName string) error {
	conn, ch, err := r.connect(rabbitURL)
	defer func() {
		conn.Close()
		ch.Close()
	}()

	// Ensure exchange exists
	err = ch.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	queueName := "events." + consumerName

	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		consumerName,
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("consumer %s waiting for messages", consumerName)

	for msg := range msgs {
		log.Printf("[%s] received: %s", consumerName, msg.Body)
		msg.Ack(false)

		// _ = msg.Nack(false, true) // Requeue
	}

	return nil
}
