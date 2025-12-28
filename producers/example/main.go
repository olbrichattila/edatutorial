package main

import event "eda.event"

func main() {
	eventManager := event.New()

	eventManager.Publish("qname", []byte("Hello"))
}
