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
	topic    = "payment-failed"
	consumer = "order-cancellation-email-service"
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

	if err := eventManager.Consume(topic, consumer, handleSendCancellationEmail(logger, actionStoreRepository)); err != nil {
		logger.Error(fmt.Sprintf("send cancel email consumer error: %v", err))
	}
}

func handleSendCancellationEmail(
	logger loggerContracts.Logger,
	actionStoreRepository actionStoreContracts.ActionStore,
) func(evt contracts.EventManager, msg []byte) error {
	return func(evt contracts.EventManager, msg []byte) error {
		paymentFailedAction, err := actions.FromJSON[actions.PaymentFailedActon](msg)
		if err != nil {
			logger.Error("order unmarshal error: " + err.Error())
			return err
		}

		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		logger.InfoWithAction(paymentFailedAction.MetaData, log)

		isDuplicate, _, err := actionStoreRepository.IsDuplicateAction(consumer, paymentFailedAction.MetaData)
		if err != nil {
			logger.ErrorWithAction(paymentFailedAction.MetaData, "deduplicate action: "+err.Error())
			return err
		}

		if isDuplicate {
			return nil
		}

		emailBody := fmt.Sprintf(`<html>
			<body>
				<h2>Hello</h2>
				<p>We are regret to inform, that your payment is failed therefore we had to cancel your order!</p>
				<p>Your order reference is: %s</p>
				<p>Please try again or contact support</p>
			</body>
		</html>`,
			paymentFailedAction.Payload.OrderUUID,
		)

		err = notification.SendEmail(paymentFailedAction.Payload.Email, "Order cancellation", emailBody)
		if err != nil {
			logger.ErrorWithAction(paymentFailedAction.MetaData, "send email error: "+err.Error())
			return err
		}

		return actionStoreRepository.SetDuplicateAction(consumer, paymentFailedAction.MetaData, "")
	}
}
