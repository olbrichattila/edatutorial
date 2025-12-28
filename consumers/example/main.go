package main

import event "eda.event"

func main() {
	eventManager := event.New()

	eventManager.Consume("consumer1")
}
