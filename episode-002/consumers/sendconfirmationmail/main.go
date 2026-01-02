package main

import (
	"fmt"
	"os"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	"github.com/olbrichattila/edatutorial/shared/logger/eventlogger"
	"github.com/olbrichattila/edatutorial/shared/notification"
)

const (
	topic    = "paymentdone"
	consumer = "sendconfirmemail"
)

func main() {
	eventManager, err := event.New()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	logger := eventlogger.New(eventManager)

	eventManager.Consume(topic, consumer, func(evt contracts.EventManager, msg []byte) error {
		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		logger.Info(log)

		orderSent, err := actions.FromJSON[actions.OrderStoredAction](msg)
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
			orderSent.Payload.ID,
		)

		err = notification.SendEmail(orderSent.Payload.Email, "Order Confirmation", emailBody)
		if err != nil {
			logger.Error("send email error: " + err.Error())
			return err
		}

		return nil
	})
}
