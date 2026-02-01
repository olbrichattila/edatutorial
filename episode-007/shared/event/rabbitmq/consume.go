package rabbitmq

import (
	"log"
	"time"

	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Episode 007, use a single dead letter queue for all messages.
	// For deciding where it comes from you could use in your dead letter consumer
	// 		xDeathHeader := msg.Headers["x-death"].([]interface{})
	// 		firstDeath := xDeathHeader[0].(amqp.Table)
	// 		originalExchange := firstDeath["exchange"].(string)
	// 		originalRoutingKey := firstDeath["routing-keys"].([]interface{})[0].(string)
	// The x-death header tells you the retry count:
	//	 "x-death": [
	//	 {
	//		"count": 2,
	//	 	"exchange": "main-ex",
	//	 	"queue": "main-queue",
	//	 	"reason": "rejected",
	//	 	"time": "2026-02-01T16:00:00"
	//	 }
	//	]

	ttlInSeconds = 30

	dlxName = "dead-letters.dlx"
	dlqName = "dead-letters"
)

func (r *rb) Consume(topic, consumerName string, callback contracts.CallbackFn) error {
	conn, ch, err := r.connect()
	if err != nil {
		return err
	}
	defer func() {
		if ch != nil {
			if closeErr := ch.Close(); closeErr != nil {
				log.Printf("Error closing channel: %v", closeErr)
			}
		}

		if conn != nil {
			if closeErr := conn.Close(); closeErr != nil {
				log.Printf("Error closing connection: %v", closeErr)
			}
		}
	}()

	// Episode 007 Declare DLQ Exchange and DLQ
	// Create Dead Letter Exchange and DLQ Queue
	// If you would like to handle all dead letters in a separate DLQ
	// dlxName := topic + "." + consumerName + ".dlx"
	// dlqName := topic + consumerName + ".dlq"

	err = r.declareOrCreateExchange(ch, dlxName)
	if err != nil {
		return err
	}

	dlq, err := r.declareQueueIfNotDeclared(ch, dlqName, "dlq")
	if err != nil {
		return err
	}

	// Bind them together
	err = r.bindQueueIfNotDoneAlready(ch, dlxName, dlq.Name)
	if err != nil {
		return err
	}

	// Declare normal data flow exchange and queue
	err = r.declareOrCreateExchange(ch, topic)
	if err != nil {
		return err
	}

	// Updated in Episode 007, use declareQueueWithDLQIfNotDeclared to setup DLX exchange
	q, err := r.declareQueueWithDLQIfNotDeclared(ch, topic, consumerName, dlxName)
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

// Episode 007 Declare DLQ
func (r *rb) declareQueueWithDLQIfNotDeclared(ch *amqp.Channel, topic, consumerName, dlxName string) (amqp.Queue, error) {
	queueName := topic + "." + consumerName
	return ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // Exclusive
		false, // NoWait
		amqp.Table{
			"x-dead-letter-exchange": dlxName,
			"x-message-ttl":          int32((ttlInSeconds * time.Second).Milliseconds()),
		},
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
			_ = msg.Nack(false, false) // Requeue changed to false in Episode 007 to use DLQ
			continue
		}

		if ackErr := msg.Ack(false); ackErr != nil {
			log.Printf("Error acknowledging message: %v", ackErr)
		}
	}

	return nil
}
