package main

import (
	"fmt"
	"strconv"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	"github.com/olbrichattila/edatutorial/shared/notification"
)

const (
	topic    = "paymentdone"
	consumer = "senconfirmemail"

	logTopic = "logmessagecreated"
)

func main() {
	eventManager := event.New()

	eventManager.Consume(topic, consumer, func(evt contracts.EventManager, msg []byte) error {
		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		fmt.Println(log)
		evt.Publish(logTopic, []byte(log))

		orderSent, err := actions.FromJSON[actions.OrderStoredAction](msg)
		if err != nil {
			evt.Publish(logTopic, []byte("send email error: "+err.Error()))
			return err
		}

		emailBody := `<html>
			<body>
				<h2>Hello</h2>
				<p>Thank you for the order</p>
				<p>Your order reference is: ` + strconv.Itoa(int(orderSent.Payload.ID)) + `
			</body>
		</html>`

		err = notification.SendEmail(orderSent.Payload.Email, "Order Confirmation", emailBody)
		if err != nil {
			evt.Publish(logTopic, []byte("send email error: "+err.Error()))
			return err
		}

		return nil
	})
}
