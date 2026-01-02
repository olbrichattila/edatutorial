package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	"github.com/olbrichattila/edatutorial/shared/logger/eventlogger"
	"producer.example/internal/repositories/invoice"
)

const (
	topic    = "paymentdone"
	consumer = "createinvoice"

	invoiceCreatedTopic = "invoicecreated"
)

func main() {
	eventManager := event.New()
	logger := eventlogger.New(eventManager)

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
		logger.Info(log)

		orderSent, err := actions.FromJSON[actions.OrderStoredAction](msg)
		if err != nil {
			logger.Error("cannot create invoice: " + err.Error())
			return err
		}

		invoiceId, err := invoiceRepository.CreateInvoice(orderSent.Payload.ID)
		if err != nil {
			logger.Error("cannot create invoice: " + err.Error())
			return err
		}

		invoiceCreatedAction := actions.New(actions.InvoiceCreatedAction{ID: invoiceId})
		invoicePayload, err := json.Marshal(invoiceCreatedAction)

		return evt.Publish(invoiceCreatedTopic, invoicePayload)
	})
}
