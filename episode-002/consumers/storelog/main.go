package main

import (
	"fmt"
	"os"

	"github.com/olbrichattila/edatutorial/shared/actions"
	"github.com/olbrichattila/edatutorial/shared/dbexecutor"
	"github.com/olbrichattila/edatutorial/shared/event"
	"github.com/olbrichattila/edatutorial/shared/event/contracts"
	"producer.example/internal/repositories/logger"
)

const (
	topic    = "logmessagecreated"
	consumer = "storelog"
)

func main() {
	eventManager := event.New()

	db, err := dbexecutor.ConnectToDB()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logRepository := logger.New(db)

	eventManager.Consume(topic, consumer, func(evt contracts.EventManager, msg []byte) error {
		fmt.Println(string(msg))

		logActionEnvelope, err := actions.FromJSON[actions.LogAction](msg)
		if err != nil {
			return err
		}

		err = logRepository.Save(
			string(logActionEnvelope.Payload.LogType),
			logActionEnvelope.ActionID,
			logActionEnvelope.Payload.Message,
		)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})
}
