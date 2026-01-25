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
	topic    = "out-of-stock"
	consumer = "order-out-of-stock-email-service"
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

	if err := eventManager.Consume(topic, consumer, handleOutOfStockEmail(logger, actionStoreRepository)); err != nil {
		logger.Error(fmt.Sprintf("send out of stock email consumer error: %v", err))
	}
}

func handleOutOfStockEmail(
	logger loggerContracts.Logger,
	actionStoreRepository actionStoreContracts.ActionStore,
) func(evt contracts.EventManager, msg []byte) error {
	return func(evt contracts.EventManager, msg []byte) error {
		outOfStockAction, err := actions.FromJSON[actions.OutOfStockAction](msg)
		if err != nil {
			logger.Error("order unmarshal error: " + err.Error())
			return err
		}

		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		logger.InfoWithAction(outOfStockAction.MetaData, log)

		isDuplicate, _, err := actionStoreRepository.IsDuplicateAction(consumer, outOfStockAction.MetaData)
		if err != nil {
			logger.ErrorWithAction(outOfStockAction.MetaData, "deduplicate action: "+err.Error())
			return err
		}

		if isDuplicate {
			return nil
		}

		emailBody := fmt.Sprintf(`<html>
			<body>
				<h2>Hello</h2>
				<p>We are regret to inform, that the stock was low when you placed the order, and someone else ordered the item before you!</p>
				<p>Your order reference is: %s</p>
				<p>If you had multiple items on your order, check if any other items are still available and try again</p>
			</body>
		</html>`,
			outOfStockAction.Payload.OrderUUID,
		)

		err = notification.SendEmail(outOfStockAction.Payload.Email, "Order out of stock", emailBody)
		if err != nil {
			logger.ErrorWithAction(outOfStockAction.MetaData, "send email error: "+err.Error())
			return err
		}

		return actionStoreRepository.SetDuplicateAction(consumer, outOfStockAction.MetaData, "")
	}
}
