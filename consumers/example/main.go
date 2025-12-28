package main

import event "github.com/olbrichattila/edatutorial"

func main() {
	eventManager := event.New()

	eventManager.Consume("consumer1")
}
