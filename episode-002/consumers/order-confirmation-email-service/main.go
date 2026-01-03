package main

import (
	"fmt"
	"os"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	loggerContracts "github.com/olbrichattila/edatutorial/shared/logger/contracts"
	"github.com/olbrichattila/edatutorial/shared/logger/eventlogger"
	"github.com/olbrichattila/edatutorial/shared/notification"
)

const (
	topic    = "payment-succeeded"
	consumer = "order-confirmation-email-service"
)

func main() {
	eventManager, err := event.New()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	logger := eventlogger.New(eventManager)

	if err := eventManager.Consume(topic, consumer, handleSendConfirmationEmail(logger)); err != nil {
		logger.Error(fmt.Sprintf("send confirmation email consumer error: %v", err))
	}
}

func handleSendConfirmationEmail(logger loggerContracts.Logger) func(evt contracts.EventManager, msg []byte) error {
	return func(evt contracts.EventManager, msg []byte) error {
		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		logger.Info(log)

		OrderCreated, err := actions.FromJSON[actions.OrderPersistedAction](msg)
		if err != nil {
			logger.Error("send email error: " + err.Error())
			return err
		}

		emailBody := fmt.Sprintf(`<html>
			<body>
				<h2>Hello</h2>
				<p>Thank you for the order</p>
				<p>Your order reference is: %d</p>
			</body>
		</html>`,
			OrderCreated.Payload.ID,
		)

		err = notification.SendEmail(OrderCreated.Payload.Email, "Order Confirmation", emailBody)
		if err != nil {
			logger.Error("send email error: " + err.Error())
			return err
		}

		return nil
	}
}
