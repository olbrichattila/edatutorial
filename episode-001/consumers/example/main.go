package main

import (
	"fmt"

	event "github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
)

const (
	topic    = "order"
	consumer = "store"
)

func main() {
	eventManager := event.New()

	eventManager.Consume(topic, consumer, func(evt contracts.EventManager, msg []byte) error {
		fmt.Printf("message received: %s\n", string(msg))

		// evt.Publish("topic2", []byte("New message"))

		return nil
	})
}
