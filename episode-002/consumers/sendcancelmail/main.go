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
	topic    = "paymentfailed"
	consumer = "sendcancelemail"
)

func main() {
	eventManager, err := event.New()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	logger := eventlogger.New(eventManager)

	if err := eventManager.Consume(topic, consumer, handleSendCancellationEmail(logger)); err != nil {
		logger.Error(fmt.Sprintf("send cancel email consumer error: %v", err))
	}
}

func handleSendCancellationEmail(logger loggerContracts.Logger) func(evt contracts.EventManager, msg []byte) error {
	return func(evt contracts.EventManager, msg []byte) error {
		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		logger.Info(log)

		orderSent, err := actions.FromJSON[actions.OrderStoredAction](msg)
		if err != nil {
			logger.Error("order unmarshal error: " + err.Error())
			return err
		}

		emailBody := fmt.Sprintf(`<html>
			<body>
				<h2>Hello</h2>
				<p>We are regret to inform, that your payment is failed therefore we had to cancel your order!</p>
				<p>Your order reference is: %d</p>
				<p>Please try again or contact support</p>
			</body>
		</html>`,
			orderSent.Payload.ID,
		)

		err = notification.SendEmail(orderSent.Payload.Email, "Order cancellation", emailBody)
		if err != nil {
			logger.Error("send email error: " + err.Error())
			return err
		}

		return nil
	}
}
