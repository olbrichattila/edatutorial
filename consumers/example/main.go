package main

import event "github.com/olbrichattila/edatutorial/shared/event"

func main() {
	eventManager := event.New()

	eventManager.Consume("consumer1")
}
