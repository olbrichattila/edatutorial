package main

import (
	"errors"
	"fmt"
	"os"
	"time"

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
	topicOutOfStock     = "out-of-stock"

	retryCount          = 5
	backoffMilliseconds = 50
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

		retried := 0

		// Episode 005 Refactored to handle race conditions, optimistic way
		for {
			_, err := orderRepository.Save(orderCreatedAction.Payload)
			if err == nil {
				break
			}

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

			// Episode 005 Ran out of stock, this is not recoverable, do not retry, publish action to send out of stock email
			if errors.Is(err, order.ErrNotEnoughStock) {
				outOfStockAction := actions.NewFromParent(topicOrderPersisted, orderCreatedAction, 0, actions.OutOfStockAction{
					OrderUUID: orderCreatedAction.Payload.UUID,
					Email:     orderCreatedAction.Payload.Email,
				})

				orderJson, err := outOfStockAction.ToJSON()
				if err != nil {
					logger.ErrorWithAction(orderCreatedAction.MetaData, fmt.Sprintf("cannot create json for stored order: %v", err))
					return err
				}

				logger.Info(fmt.Sprintf("insufficient stock: Order: %s", orderCreatedAction.Payload.UUID))

				return evt.Publish(topicOutOfStock, orderJson)
			}

			// Episode 005 Race condition during stock update, multiple retry
			if errors.Is(err, order.ErrRaceConditionDetected) {
				retried++

				// Retry {retryCount} times
				if retried == retryCount {
					msg := "order service, race condition retry count reached"
					logger.ErrorWithAction(orderCreatedAction.MetaData, msg)
					return errors.New(msg)
				}

				// Backoff time is 50 milliseconds, note consider progressive retry backoff times
				time.Sleep(backoffMilliseconds * time.Millisecond)
				continue
			}

			//
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
	orderPersistedAction := actions.NewFromParent(topicOrderPersisted, parent, 0, actions.OrderPersistedAction{
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
