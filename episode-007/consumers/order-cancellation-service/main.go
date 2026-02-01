package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	loggerContracts "github.com/olbrichattila/edatutorial/shared/logger/contracts"
	"github.com/olbrichattila/edatutorial/shared/logger/eventlogger"
	repositoryContracts "producer.example/internal/contracts"
	"producer.example/internal/repositories/order"
)

const (
	topic    = "payment-failed"
	consumer = "order-cancellation-service"

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

	if err := eventManager.Consume(topic, consumer, handleCancelOrder(logger, orderRepository)); err != nil {
		logger.Error(fmt.Sprintf("cancer order consumer error: %v", err))
	}
}

func handleCancelOrder(logger loggerContracts.Logger, orderRepository repositoryContracts.OrderRepository) func(evt contracts.EventManager, msg []byte) error {
	return func(evt contracts.EventManager, msg []byte) error {
		paymentFailedAction, err := actions.FromJSON[actions.PaymentFailedActon](msg)
		if err != nil {
			logger.Error("cannot cancel order: " + err.Error())
			return err
		}

		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		logger.InfoWithAction(paymentFailedAction.MetaData, log)

		// Episode 005 Race condition
		retries := 0
		for {
			err := orderRepository.Cancel(paymentFailedAction.Payload.OrderUUID)
			if err == nil {
				// Normal flow
				break
			}

			// If race condition detected
			if errors.Is(err, order.ErrRaceConditionDetected) {
				retries++

				// Retry {retryCount} times
				if retries == retryCount {
					msg := "order cancellation, retry count reached"
					logger.ErrorWithAction(paymentFailedAction.MetaData, msg)
					return errors.New(msg)
				}

				// Wait for the other operation to finish, Note, progressive timing also can be implemented
				time.Sleep(backoffMilliseconds * time.Millisecond)
				continue
			}

			// Non race condition handler, no retry
			return err
		}

		return nil
	}
}
