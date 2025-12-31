package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	"producer.example/internal/dto"
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

		var order dto.Order
		err := json.Unmarshal(msg, &order)
		if err != nil {
			return evt.Publish(logTopic, []byte(err.Error()))
		}

		orderId, err := orderRepository.Save(order)
		if err != nil {
			return evt.Publish(logTopic, []byte(err.Error()))
		}

		orderAction := dto.OrderAction{ID: orderId, Email: order.Email}
		orderJson, err := json.Marshal(orderAction)
		if err != nil {
			return evt.Publish(logTopic, []byte(err.Error()))
		}

		return evt.Publish(topicOrderStored, orderJson)
	})
}
