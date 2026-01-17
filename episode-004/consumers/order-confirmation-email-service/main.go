package main

import (
	"fmt"
	"os"

	"github.com/olbrichattila/edatutorial/shared/actions"
	actionStoreContracts "github.com/olbrichattila/edatutorial/shared/actionstore/contracts"
	"github.com/olbrichattila/edatutorial/shared/actionstore/mysqlactionstore"
	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
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

	db, err := dbexecutor.ConnectToDB()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	actionStoreRepository, err := mysqlactionstore.New(db)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := eventManager.Consume(topic, consumer, handleSendConfirmationEmail(logger, actionStoreRepository)); err != nil {
		logger.Error(fmt.Sprintf("send confirmation email consumer error: %v", err))
	}
}

func handleSendConfirmationEmail(
	logger loggerContracts.Logger,
	actionStoreRepository actionStoreContracts.ActionStore,

) func(evt contracts.EventManager, msg []byte) error {
	return func(evt contracts.EventManager, msg []byte) error {
		paymentSucceededAction, err := actions.FromJSON[actions.PaymentSucceededAction](msg)
		if err != nil {
			logger.Error("send email error: " + err.Error())
			return err
		}

		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		logger.InfoWithAction(paymentSucceededAction.MetaData, log)

		// Episode 004 Add deduplication store
		isDuplicate, err := actionStoreRepository.IsDuplicateAction(paymentSucceededAction.MetaData)
		if err != nil {
			logger.ErrorWithAction(paymentSucceededAction.MetaData, "deduplicate action: "+err.Error())
			return err
		}

		if isDuplicate {
			return nil
		}

		emailBody := fmt.Sprintf(`<html>
			<body>
				<h2>Hello</h2>
				<p>Thank you for the order</p>
				<p>Your order reference is: %s</p>
			</body>
		</html>`,
			paymentSucceededAction.Payload.OrderUUID,
		)

		err = notification.SendEmail(paymentSucceededAction.Payload.Email, "Order Confirmation", emailBody)
		if err != nil {
			logger.ErrorWithAction(paymentSucceededAction.MetaData, "send email error: "+err.Error())
			return err
		}

		return nil
	}
}
