package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/olbrichattila/edatutorial/shared/actions"
	actionStoreContracts "github.com/olbrichattila/edatutorial/shared/actionstore/contracts"
	"github.com/olbrichattila/edatutorial/shared/actionstore/mysqlactionstore"
	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	loggerContracts "github.com/olbrichattila/edatutorial/shared/logger/contracts"
	"github.com/olbrichattila/edatutorial/shared/logger/eventlogger"
)

const (
	topic    = "order-persisted"
	consumer = "payment-service"

	paymentSucceededTopic = "payment-succeeded"
	paymentFailedTopic    = "payment-failed"

	successMetaData = "succeeded"
	failedMetaData  = "failed"
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

	if err := eventManager.Consume(topic, consumer, handlePayment(logger, actionStoreRepository)); err != nil {
		logger.Error(fmt.Sprintf("payment consumer error: %v", err))
	}
}

func handlePayment(
	logger loggerContracts.Logger,
	actionStoreRepository actionStoreContracts.ActionStore,
) func(evt contracts.EventManager, msg []byte) error {
	return func(evt contracts.EventManager, msg []byte) error {
		orderPersistedAction, err := actions.FromJSON[actions.OrderPersistedAction](msg)
		if err != nil {
			logger.Error("cannot create invoice: " + err.Error())
			return err
		}

		log := fmt.Sprintf("topic: %s, consumer: %s, message %s\n", topic, consumer, string(msg))
		logger.InfoWithAction(orderPersistedAction.MetaData, log)

		// Episode 004 Add deduplication store
		isDuplicate, duplicationMetadata, err := actionStoreRepository.IsDuplicateAction(consumer, orderPersistedAction.MetaData)
		if err != nil {
			logger.ErrorWithAction(orderPersistedAction.MetaData, "deduplicate action: "+err.Error())
			return err
		}

		paymentSuccess := false
		if isDuplicate {
			fmt.Println("duplicate", successMetaData, duplicationMetadata == successMetaData)
			paymentSuccess = duplicationMetadata == successMetaData
		} else {
			paymentSuccess = randomPaymentSuccess()
			fmt.Println("new", paymentSuccess)
			time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
		}

		if paymentSuccess {
			paymentAction := actions.NewFromParent(paymentSucceededTopic, orderPersistedAction, 0, actions.PaymentSucceededAction{
				OrderUUID: orderPersistedAction.Payload.OrderUUID,
				Email:     orderPersistedAction.Payload.Email,
			})
			asJSON, err := paymentAction.ToJSON()
			if err != nil {
				return err
			}
			err = evt.Publish(paymentSucceededTopic, asJSON)
			if err != nil {
				return err
			}

			return actionStoreRepository.SetDuplicateAction(consumer, orderPersistedAction.MetaData, successMetaData)
		}

		paymentAction := actions.NewFromParent(paymentFailedTopic, orderPersistedAction, 0, actions.PaymentFailedActon{
			OrderUUID: orderPersistedAction.Payload.OrderUUID,
			Email:     orderPersistedAction.Payload.Email,
		})
		asJSON, err := paymentAction.ToJSON()
		if err != nil {
			return err
		}

		err = evt.Publish(paymentFailedTopic, asJSON)
		if err != nil {
			return err
		}

		return actionStoreRepository.SetDuplicateAction(consumer, orderPersistedAction.MetaData, failedMetaData)
	}
}

func randomPaymentSuccess() bool {
	if rand.Intn(10) > 7 {
		return false
	}

	return true
}
