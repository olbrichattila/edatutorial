package main

import (
	"fmt"

	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
)

const (
	topic    = "paymentdone"
	consumer = "senconfirmemail"

	logTopic = "logmessagecreated"
)

func main() {
	eventManager := event.New()

	eventManager.Consume(topic, consumer, func(evt contracts.EventManager, msg []byte) error {
		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		fmt.Println(log)
		evt.Publish(logTopic, []byte(log))

		return nil
	})
}
