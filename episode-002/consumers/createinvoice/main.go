package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	"producer.example/internal/repositories/invoice"
)

const (
	topic    = "paymentdone"
	consumer = "createinvoice"

	logTopic = "logmessagecreated"
)

type orderAction struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func main() {
	eventManager := event.New()

	db, err := dbexecutor.ConnectToDB()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	invoiceRepository := invoice.New(db)

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

		err = invoiceRepository.CreateInvoice(order.ID)
		if err != nil {
			evt.Publish(logTopic, []byte("cannot create invoice: "+err.Error()))
			return err
		}

		return nil
	})
}
