package main

import event "github.com/olbrichattila/edatutorial/tree/main/shared/event/contracts"

func main() {
	eventManager := event.New()

	eventManager.Publish("qname", []byte("Hello"))
}
