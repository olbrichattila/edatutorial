package main

import (
	"fmt"

	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
)

const (
	topic    = "logmessagecreated"
	consumer = "storelog"
)

func main() {
	eventManager := event.New()

	eventManager.Consume(topic, consumer, func(evt contracts.EventManager, msg []byte) error {
		fmt.Println(string(msg))

		return nil
	})
}
