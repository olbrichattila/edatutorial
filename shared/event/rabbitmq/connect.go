package rabbitmq

import (
	"github.com/olbrichattila/edatutorial/shared/event/contracts"

	amqp "github.com/rabbitmq/amqp091-go"
)

func New() contracts.EventManager {
	return &rb{}
}

type rb struct {
}

func (r *rb) connect(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	return conn, ch, nil
}
