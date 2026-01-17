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
	orderContracts "producer.example/internal/contracts"
	"producer.example/internal/repositories/order"
)

const (
	topic    = "order-created"
	consumer = "order-service"

	topicOrderPersisted = "order-persisted"
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

	orderRepository := order.New(db)

	if err := eventManager.Consume(topic, consumer, handleStoreOrder(logger, orderRepository)); err != nil {
		logger.Error(fmt.Sprintf("store order consumer error: %v", err))
	}
}

func handleStoreOrder(logger loggerContracts.Logger, orderRepository orderContracts.OrderRepository) func(evt contracts.EventManager, msg []byte) error {
	return func(evt contracts.EventManager, msg []byte) error {
		orderCreatedAction, err := actions.FromJSON[actions.OrderCreatedAction](msg)
		if err != nil {
			logger.Error(fmt.Sprintf("cannot get sent order: %v", err))
			return err
		}

		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		logger.InfoWithAction(orderCreatedAction.MetaData, log)

		_, err = orderRepository.Save(orderCreatedAction.Payload)
		if err != nil {
			if isDuplicateKeyError(err) {
				// Episode 004 Re throw event
				return publishOrderPersistedAction(
					orderCreatedAction,
					logger,
					evt,
					orderCreatedAction.Payload.UUID, // Episode 004 de-dupe invoice as well
					orderCreatedAction.Payload.Email,
				)
			}

			logger.ErrorWithAction(orderCreatedAction.MetaData, fmt.Sprintf("cannot save order: %v", err))
			return err
		}

		return publishOrderPersistedAction(
			orderCreatedAction,
			logger,
			evt,
			orderCreatedAction.Payload.UUID, // Episode 004 de-dupe invoice as well
			orderCreatedAction.Payload.Email,
		)
	}
}

// Episode 004 moved out as needs to be triggered if we skip due to de-dupe
func publishOrderPersistedAction(
	parent actions.Envelope[actions.OrderCreatedAction],
	logger loggerContracts.Logger,
	evt contracts.EventManager,
	orderUUID string,
	email string,
) error {
	orderPersistedAction := actions.NewFromParent(parent, 0, actions.OrderPersistedAction{
		OrderUUID: orderUUID, // Episode 004
		Email:     email,
	})

	orderJson, err := orderPersistedAction.ToJSON()
	if err != nil {
		logger.ErrorWithAction(parent.MetaData, fmt.Sprintf("cannot create json for stored order: %v", err))
		return err
	}

	return evt.Publish(topicOrderPersisted, orderJson)
}

func isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}
