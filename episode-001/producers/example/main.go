package main

import event "github.com/olbrichattila/edatutorial/shared/event"

const topic = "order"

func main() {
	eventManager := event.New()

	eventManager.Publish(topic, []byte("Hello"))
}
