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
	eventManager, err := event.New()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	logger := eventlogger.New(eventManager)

	db, err := dbexecutor.ConnectToDB()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func() {
		if db != nil {
			if closeErr := db.Close(); closeErr != nil {
				fmt.Printf("Error closing database: %v\n", closeErr)
			}
		}
	}()

	orderRepository := order.New(db)

	eventManager.Consume(topic, consumer, func(evt contracts.EventManager, msg []byte) error {
		log := fmt.Sprintf("topic: %s, consumer: %s, message received\n", topic, consumer)
		logger.Info(log)

		envelope, err := actions.FromJSON[actions.OrderSentAction](msg)
		if err != nil {
			logger.Error(fmt.Sprintf("cannot get sent order: %v", err))
			return err
		}

		orderID, err := orderRepository.Save(envelope.Payload)
		if err != nil {
			logger.Error(fmt.Sprintf("cannot save order: %v", err))
			return err
		}

		orderAction := actions.New(actions.OrderStoredAction{ID: orderID, Email: envelope.Payload.Email})
		orderJson, err := json.Marshal(orderAction)
		if err != nil {
			logger.Error(fmt.Sprintf("cannot create json for stored order: %v", err))
			return err
		}

		return evt.Publish(topicOrderStored, orderJson)
	})
}
