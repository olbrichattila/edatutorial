package main

import event "github.com/olbrichattila/edatutorial/tree/main/shared/event/contracts"

func main() {
	eventManager := event.New()

	eventManager.Consume("consumer1")
}
