package main

import (
	"fmt"
	"os"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	"producer.example/internal/repositories/order"
)

const (
	topic    = "paymentfailed"
	consumer = "cancelorder"

	logTopic = "logmessagecreated"
)

func main() {
	eventManager := event.New()

	db, err := dbexecutor.ConnectToDB()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	orderRepository := order.New(db)

	eventManager.Consume(topic, consumer, func(evt contracts.EventManager, msg []byte) error {
		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		fmt.Println(log)
		evt.Publish(logTopic, []byte(log))

		orderSent, err := actions.FromJSON[actions.OrderStoredAction](msg)
		if err != nil {
			evt.Publish(logTopic, []byte("cannot cancel order: "+err.Error()))
			return err
		}

		err = orderRepository.Cancel(orderSent.Payload.ID)
		if err != nil {
			evt.Publish(logTopic, []byte("cannot cancel order: "+err.Error()))
			return err
		}

		return nil
	})
}
