package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	loggerContracts "github.com/olbrichattila/edatutorial/shared/logger/contracts"
	"github.com/olbrichattila/edatutorial/shared/logger/eventlogger"
)

const (
	topic    = "order-persisted"
	consumer = "payment-service"

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
		orderPersistedAction, err := actions.FromJSON[actions.OrderPersistedAction](msg)
		if err != nil {
			logger.Error("cannot create invoice: " + err.Error())
			return err
		}

		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		logger.InfoWithAction(orderPersistedAction.MetaData, log)

		// Random wait, emulate user pays
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)

		// TODO this forwards the same action, it should be paymentSucceeded or PaymentFailed action, though the structure is the same

		if paymentSuccess() {
			paymentAction := actions.NewFromParent(orderPersistedAction, 0, actions.PaymentSucceededAction{
				OrderUUID: orderPersistedAction.Payload.OrderUUID,
				Email:     orderPersistedAction.Payload.Email,
			})
			asJSON, err := paymentAction.ToJSON()
			if err != nil {
				return err
			}
			return evt.Publish(paymentSucceededTopic, asJSON)
		}

		paymentAction := actions.NewFromParent(orderPersistedAction, 0, actions.PaymentFailedActon{
			OrderUUID: orderPersistedAction.Payload.OrderUUID,
			Email:     orderPersistedAction.Payload.Email,
		})
		asJSON, err := paymentAction.ToJSON()
		if err != nil {
			return err
		}

		return evt.Publish(paymentFailedTopic, asJSON)
	}
}

func paymentSuccess() bool {
	if rand.Int63n(10) > 7 {
		return false
	}

	return true
}
