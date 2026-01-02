package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	"github.com/olbrichattila/edatutorial/shared/logger/eventlogger"
	"producer.example/internal/repositories/order"
)

const (
	topic    = "order"
	consumer = "store"

	topicOrderStored = "orderstored"
)

func main() {
	eventManager := event.New()
	logger := eventlogger.New(eventManager)

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
		logger.Info(log)

		envelope, err := actions.FromJSON[actions.OrderSentAction](msg)
		if err != nil {
			logger.Error("cannot get sent order =: " + err.Error())
			return err
		}

		orderId, err := orderRepository.Save(envelope.Payload)
		if err != nil {
			logger.Error("cannot save order =: " + err.Error())
			return err
		}

		orderAction := actions.New(actions.OrderStoredAction{ID: orderId, Email: envelope.Payload.Email})
		orderJson, err := json.Marshal(orderAction)
		if err != nil {
			logger.Error("cannot create order action: " + err.Error())
			return err
		}

		return evt.Publish(topicOrderStored, orderJson)
	})
}
