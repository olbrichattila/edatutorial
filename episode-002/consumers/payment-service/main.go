package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	loggerContracts "github.com/olbrichattila/edatutorial/shared/logger/contracts"
	"github.com/olbrichattila/edatutorial/shared/logger/eventlogger"
)

const (
	topic    = "order-persisted"
	consumer = "order-service"

	paymentSucceededTopic = "payment-succeeded"
	paymentFailedTopic    = "payment-failed"
)

func main() {
	eventManager, err := event.New()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	logger := eventlogger.New(eventManager)

	if err := eventManager.Consume(topic, consumer, handlePayment(logger)); err != nil {
		logger.Error(fmt.Sprintf("payment consumer error: %v", err))
	}
}

func handlePayment(logger loggerContracts.Logger) func(evt contracts.EventManager, msg []byte) error {
	return func(evt contracts.EventManager, msg []byte) error {
		log := fmt.Sprintf("topic: %s, consumer: %s, message received", topic, consumer)
		logger.Info(log)

		// Random wait, emulate user pays
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)

		if paymentSuccess() {
			return evt.Publish(paymentSucceededTopic, msg)
		}

		return evt.Publish(paymentFailedTopic, msg)
	}
}

func paymentSuccess() bool {
	if rand.Int63n(10) > 7 {
		return false
	}

	return true
}
