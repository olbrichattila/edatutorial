package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	"producer.example/internal/repositories/order"
)

const (
	topic    = "order"
	consumer = "store"

	topicOrderStored = "orderstored"
	logTopic         = "logmessagecreated"
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

		envelope, err := actions.FromJSON[actions.OrderSentAction](msg)
		if err != nil {
			return evt.Publish(logTopic, []byte(err.Error()))
		}

		orderId, err := orderRepository.Save(envelope.Payload)
		if err != nil {
			return evt.Publish(logTopic, []byte(err.Error()))
		}

		orderAction := actions.New(actions.OrderStoredAction{ID: orderId, Email: envelope.Payload.Email})
		orderJson, err := json.Marshal(orderAction)
		if err != nil {
			return evt.Publish(logTopic, []byte(err.Error()))
		}

		return evt.Publish(topicOrderStored, orderJson)
	})
}
