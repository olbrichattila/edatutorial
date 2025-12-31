package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
)

const (
	topic    = "orderstored"
	consumer = "payment"

	paymentDoneTopic   = "paymentdone"
	paymentFailedTopic = "paymentfailed"
	logTopic           = "logmessagecreated"
)

func main() {
	eventManager := event.New()

	eventManager.Consume(topic, consumer, func(evt contracts.EventManager, msg []byte) error {
		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		fmt.Println(log)
		evt.Publish(logTopic, []byte(log))

		// Random wait, emulate user pays
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)

		if paymentSuccess() {
			return evt.Publish(paymentDoneTopic, msg)
		}

		return evt.Publish(paymentFailedTopic, msg)
	})
}

func paymentSuccess() bool {
	if rand.Intn(10) > 7 {
		return false
	}

	return true
}
