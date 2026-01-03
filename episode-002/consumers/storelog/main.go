package main

import (
	"fmt"
	"os"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	logRepositoryContracts "producer.example/internal/contracts"
	"producer.example/internal/repositories/logger"
)

const (
	topic    = "logmessagecreated"
	consumer = "storelog"
)

func main() {
	eventManager, err := event.New()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	db, err := dbexecutor.ConnectToDB()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer func() {
		if db != nil {
			if closeErr := db.Close(); closeErr != nil {
				fmt.Printf("Error closing database: %v\n", closeErr)
			}
		}
	}()

	logRepository := logger.New(db)

	if err := eventManager.Consume(topic, consumer, handleStoreLog(logRepository)); err != nil {
		logAction := actions.New(actions.LogAction{LogType: actions.LogTypeError, Message: fmt.Sprintf("log consumer error: %v", err)})
		_ = logRepository.Save(
			string(logAction.Payload.LogType),
			logAction.ActionID,
			logAction.Payload.Message,
		)
	}
}

func handleStoreLog(logRepository logRepositoryContracts.LoggerRepository) func(evt contracts.EventManager, msg []byte) error {
	return func(evt contracts.EventManager, msg []byte) error {
		fmt.Println(string(msg))

		logActionEnvelope, err := actions.FromJSON[actions.LogAction](msg)
		if err != nil {
			return err
		}

		if err := logRepository.Save(
			string(logActionEnvelope.Payload.LogType),
			logActionEnvelope.ActionID,
			logActionEnvelope.Payload.Message,
		); err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	}
}
