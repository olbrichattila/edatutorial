package main

import (
	"fmt"
	"os"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	loggerContracts "github.com/olbrichattila/edatutorial/shared/logger/contracts"
	"github.com/olbrichattila/edatutorial/shared/logger/eventlogger"
	repositoryContracts "producer.example/internal/contracts"
	"producer.example/internal/repositories/invoice"
)

const (
	topic    = "payment-succeeded"
	consumer = "invoice-service"

	invoiceCreatedTopic = "invoice-created"
)

func main() {
	eventManager, err := event.New()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	logger := eventlogger.New(eventManager)

	db, err := dbexecutor.ConnectToDB()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func() {
		if db != nil {
			if closeErr := db.Close(); closeErr != nil {
				fmt.Printf("Error closing database: %v\n", closeErr)
			}
		}
	}()

	invoiceRepository := invoice.New(db)

	if err := eventManager.Consume(topic, consumer, handleCreateInvoice(logger, invoiceRepository)); err != nil {
		logger.Error(fmt.Sprintf("cancer invoice consumer error: %v", err))
	}

}

func handleCreateInvoice(logger loggerContracts.Logger, invoiceRepository repositoryContracts.InvoiceRepository) func(evt contracts.EventManager, msg []byte) error {
	return func(evt contracts.EventManager, msg []byte) error {
		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		logger.Info(log)

		orderPersisted, err := actions.FromJSON[actions.OrderPersistedAction](msg)
		if err != nil {
			logger.Error("cannot create invoice: " + err.Error())
			return err
		}

		invoiceID, err := invoiceRepository.CreateInvoice(orderPersisted.Payload.ID)
		if err != nil {
			logger.Error("cannot create invoice: " + err.Error())
			return err
		}

		InvoiceCreatedAction := actions.New(actions.InvoiceCreatedAction{ID: invoiceID})
		invoicePayload, err := InvoiceCreatedAction.ToJSON()
		if err != nil {
			logger.Error("cannot get invoice payload: " + err.Error())
			return err
		}

		return evt.Publish(invoiceCreatedTopic, invoicePayload)
	}
}
