package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	"github.com/olbrichattila/edatutorial/shared/notification"
)

const (
	topic    = "paymentfailed"
	consumer = "sendcancelemail"

	logTopic = "logmessagecreated"
)

type orderAction struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func main() {
	eventManager := event.New()

	eventManager.Consume(topic, consumer, func(evt contracts.EventManager, msg []byte) error {
		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		fmt.Println(log)
		evt.Publish(logTopic, []byte(log))

		var order orderAction
		err := json.Unmarshal(msg, &order)
		if err != nil {
			evt.Publish(logTopic, []byte("order unmarshal error: "+err.Error()))
			return err
		}

		emailBody := `<html>
			<body>
				<h2>Hello</h2>
				<p>We are regret to inform, that your payment is failed therefore we had to cancel your order!</p>
				<p>Your order reference is: ` + strconv.Itoa(int(order.ID)) + `
				<p>Please try again or contact support</p>
			</body>
		</html>`

		err = notification.SendEmail(order.Email, "Order cancellation", emailBody)
		if err != nil {
			evt.Publish(logTopic, []byte("send email error: "+err.Error()))
			return err
		}

		return nil
	})
}
