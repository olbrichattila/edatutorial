package main

import event "github.com/olbrichattila/edatutorial/shared/event"

func main() {
	eventManager := event.New()

	eventManager.Publish("qname", []byte("Hello"))
}
