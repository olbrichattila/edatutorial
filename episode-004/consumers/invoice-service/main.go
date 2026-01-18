package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
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
		logger.Error(fmt.Sprintf("cancel invoice consumer error: %v", err))
	}

}

func handleCreateInvoice(logger loggerContracts.Logger, invoiceRepository repositoryContracts.InvoiceRepository) func(evt contracts.EventManager, msg []byte) error {
	return func(evt contracts.EventManager, msg []byte) error {
		paymentSucceededAction, err := actions.FromJSON[actions.PaymentSucceededAction](msg)
		if err != nil {
			logger.Error("cannot create invoice: " + err.Error())
			return err
		}

		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		logger.InfoWithAction(paymentSucceededAction.MetaData, log)

		// Episode 004 Added order UUID for de-dupe key
		invoiceID, err := invoiceRepository.CreateInvoice(paymentSucceededAction.Payload.OrderUUID)
		if err != nil {
			if isDuplicateKeyError(err) {

				id, err := invoiceRepository.GetInvoiceId(paymentSucceededAction.Payload.OrderUUID)
				if err != nil {
					logger.ErrorWithAction(paymentSucceededAction.MetaData, "cannot get invoice id in deduplication: "+err.Error())
					return fmt.Errorf("cannot get invoice id %w", err)
				}
				// Episode 004 Skip if invoice already stored, Idempotency
				return publishInvoiceCreatedAction(paymentSucceededAction, logger, evt, id, paymentSucceededAction.Payload.OrderUUID)
			}

			logger.ErrorWithAction(paymentSucceededAction.MetaData, "cannot create invoice: "+err.Error())
			return err
		}

		return publishInvoiceCreatedAction(paymentSucceededAction, logger, evt, invoiceID, paymentSucceededAction.Payload.OrderUUID)
	}
}

// Episode 004 moved out as we need to trigger the action if we skip because of idempotency
func publishInvoiceCreatedAction(
	parent actions.Envelope[actions.PaymentSucceededAction],
	logger loggerContracts.Logger,
	evt contracts.EventManager,
	invoiceID int64,
	orderUUID string,
) error {
	invoiceCreatedAction := actions.NewFromParent(invoiceCreatedTopic, parent, 0, actions.InvoiceCreatedAction{ID: invoiceID})
	invoicePayload, err := invoiceCreatedAction.ToJSON()
	if err != nil {
		logger.ErrorWithAction(parent.MetaData, "cannot get invoice payload: "+err.Error())
		return err
	}

	return evt.Publish(invoiceCreatedTopic, invoicePayload)
}

func isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}
